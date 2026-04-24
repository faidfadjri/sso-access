package oauth

import (
	"akastra-access/internal/app/config"
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/repositories"
	"akastra-access/internal/infrastructure/security/jwt"
	"akastra-access/internal/infrastructure/utils/helper"
	"akastra-access/internal/infrastructure/utils/image"
	"akastra-access/internal/interface/http/request"
	"akastra-access/internal/pkg/email"
	pkgErrors "akastra-access/internal/pkg/errors"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	sessionKeyPrefix   = "session:"
	authCodeExpiration = 10 * time.Minute
	sessionExpiration  = 24 * time.Hour
)

// AuthSession holds the state persisted in Redis during the OAuth authorization flow.
type AuthSession struct {
	ClientID            string  `json:"client_id"`
	UserID              uint64  `json:"user_id"`
	RedirectURI         string  `json:"redirect_uri"`
	Scope               string  `json:"scope"`
	State               string  `json:"state"`
	CodeChallenge       *string `json:"code_challenge,omitempty"`
	CodeChallengeMethod *string `json:"code_challenge_method,omitempty"`
}

type oauthUsecase struct {
	oauthRepository   repositories.OAuthRepository
	serviceRepository repositories.ServiceRepository
	sessionRepository repositories.SessionRepository
	assignedRolesRepository repositories.AssignedRolesRepository
	userAccessRepository repositories.UserAccessRepository
}

func NewOAuthUsecase(
	oauthRepo repositories.OAuthRepository,
	serviceRepo repositories.ServiceRepository,
	sessionRepo repositories.SessionRepository,
	assignedRolesRepository repositories.AssignedRolesRepository,
	userAccessRepo repositories.UserAccessRepository,
) OAuthUsecase {
	return &oauthUsecase{
		oauthRepository:   oauthRepo,
		serviceRepository: serviceRepo,
		sessionRepository: sessionRepo,
		assignedRolesRepository: assignedRolesRepository,
		userAccessRepository: userAccessRepo,
	}
}

// Login validates credentials, creates a server-side session, and returns the session ID.
func (u *oauthUsecase) Login(req request.LoginRequest) (string, *domain.Users, error) {
	user, err := u.oauthRepository.GetUserByEmailorUsername(req.EmailOrUsername)
	if err != nil {
		return "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", nil, pkgErrors.ErrUnauthorized
	}

	// Security: Wipe password from the returned user object to prevent accidental leaks.
	user.Password = ""

	sessionID := helper.GenerateSessionID()
	key := sessionKeyPrefix + sessionID

	if err := u.sessionRepository.SetSession(key, user.UserId, sessionExpiration); err != nil {
		return "", nil, fmt.Errorf("failed to create session: %w", err)
	}

	return sessionID, user, nil
}

// Logout removes the user's server-side session.
func (u *oauthUsecase) Logout(sessionID string) error {
	return u.sessionRepository.DeleteSession(sessionKeyPrefix + sessionID)
}

// Authorize validates the OAuth request, resolves the user from session, and issues an authorization code.
func (u *oauthUsecase) Authorize(req request.AuthorizeRequest, sessionID string) (string, error) {
	// Validate client and redirect URI in a single repository call.
	if err := u.validateClientAndRedirectURI(req.ClientID, req.RedirectURI); err != nil {
		return "", err
	}

	if sessionID == "" {
		return "", pkgErrors.ErrUnauthorized
	}

	userIDStr, err := u.sessionRepository.GetSession(sessionKeyPrefix + sessionID)
	if err != nil {
		return "", pkgErrors.ErrSessionNotFound
	}
	if userIDStr == "" {
		return "", pkgErrors.ErrSessionExpired
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return "", fmt.Errorf("failed to parse user id from session: %w", err)
	}

	authSession := AuthSession{
		ClientID:            req.ClientID,
		UserID:              userID,
		RedirectURI:         req.RedirectURI,
		Scope:               req.Scope,
		State:               req.State,
		CodeChallenge:       req.CodeChallenge,
		CodeChallengeMethod: req.CodeChallengeMethod,
	}

	sessionJSON, err := json.Marshal(authSession)
	if err != nil {
		return "", fmt.Errorf("failed to marshal auth session: %w", err)
	}

	authCode := helper.GenerateAuthCode()

	// Persist the auth code with a short TTL.
	// Note: We use SetSession explicitly instead of SaveAuthCode because we need to store the full JSON object,
	// not just the session ID, and we want a specific expiration time for the auth code flow.
	authCodeKey := u.sessionRepository.GetAuthCodePrefix() + authCode
	if err := u.sessionRepository.SetSession(authCodeKey, sessionJSON, authCodeExpiration); err != nil {
		return "", fmt.Errorf("failed to save auth code: %w", err)
	}

	return authCode, nil
}

// TokenExchange validates an authorization code and issues access + refresh tokens.
func (u *oauthUsecase) TokenExchange(req request.ExchangeAuthCodeRequest) (string, string, error) {
	// 1. Fetch Client (Service)
	service, err := u.serviceRepository.GetClientByClientId(req.ClientID)
	if err != nil {
		return "", "", fmt.Errorf("invalid client_id: %w", err)
	}
	if service == nil {
		return "", "", pkgErrors.ErrInvalidClientID
	}

	// 2. Fetch Auth Code Session
	sessionJSON, err := u.sessionRepository.GetAuthCode(req.Code)
	if err != nil {
		return "", "", pkgErrors.ErrInvalidAuthCode
	}
	if sessionJSON == "" { // Redis key might exist but be empty, or GetAuthCode semantics handle generic error
		return "", "", pkgErrors.ErrInvalidAuthCode
	}

	var authSession AuthSession
	if err := json.Unmarshal([]byte(sessionJSON), &authSession); err != nil {
		return "", "", fmt.Errorf("failed to parse auth session: %w", err)
	}

	// 3. Validate Request against Session
	if authSession.ClientID != req.ClientID {
		return "", "", pkgErrors.ErrInvalidClientID
	}
	if authSession.RedirectURI != req.RedirectURI {
		return "", "", pkgErrors.ErrInvalidRedirectURI
	}

	// 4. Validate Credentials (PKCE or Client Secret)
	if err := u.validateClientCredentials(req, authSession, service); err != nil {
		return "", "", err
	}

	// 5. Get User with Service Access
	userHasAccess, _ := u.userAccessRepository.IsUserAccessExist(authSession.UserID, service.ServiceId)
	if !userHasAccess {
		return "", "", pkgErrors.ErrUserNotHaveAccess
	}

	userWithService, err := u.oauthRepository.GetUserByIDWithService(authSession.UserID, service.ServiceId)
	if err != nil {
		return "", "", fmt.Errorf("user not found or access denied: %w", err)
	}

	// 6. Generate Tokens
	idpRole, err := u.getIDPRole(userWithService.UserId)
	if err != nil {
		return "", "", fmt.Errorf("failed to get IDP role: %w", err)
	}

	userWrapper := &tokenUserWrapper{
		TokenUser: userWithService,
		idpRole:   idpRole,
	}

	accessToken, err := jwt.GenerateAccessToken(userWrapper)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := jwt.GenerateRefreshToken(userWithService)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// 7. Consume Auth Code (Replay Protection)
	if err := u.sessionRepository.DeleteAuthCode(req.Code); err != nil {
		log.Printf("failed to delete auth code (replay protection): %v", err)
	}

	return accessToken, refreshToken, nil
}

// RefreshToken validates a refresh token and issues a new token pair.
func (u *oauthUsecase) RefreshToken(req request.RefreshTokenRequest) (string, string, error) {
	claims, err := jwt.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %w", err)
	}

	user, err := u.oauthRepository.GetUserByID(claims.UserID)
	if err != nil {
		return "", "", err
	}
	if user == nil {
		return "", "", pkgErrors.ErrNotFound
	}

	idpRole, err := u.getIDPRole(user.UserId)
	if err != nil {
		return "", "", fmt.Errorf("failed to get IDP role: %w", err)
	}
	
	userWrapper := &tokenUserWrapper{
		TokenUser: user,
		idpRole:   idpRole,
	}

	accessToken, err := jwt.GenerateAccessToken(userWrapper)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := jwt.GenerateRefreshToken(user)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// UpdateAccount applies profile changes and optionally replaces the user's photo.
func (u *oauthUsecase) UpdateAccount(userID uint64, req request.UpdateAccountRequest, file multipart.File, header *multipart.FileHeader) (string, string, error) {
	user, err := u.oauthRepository.GetUserByID(userID)
	if err != nil {
		return "", "", err
	}
	if user == nil {
		return "", "", pkgErrors.ErrNotFound
	}

	user.FullName = req.FullName
	user.Email = req.Email
	user.Username = req.Username
	user.Phone = &req.Phone

	if req.Password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return "", "", fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = string(hashed)
	}

	if file != nil {
		if err := u.replaceProfilePhoto(user, file, header); err != nil {
			return "", "", err
		}
	}

	if err := u.oauthRepository.UpdateUser(user); err != nil {
		return "", "", err
	}

	idpRole, err := u.getIDPRole(user.UserId)
	if err != nil {
		return "", "", fmt.Errorf("failed to get IDP role: %w", err)
	}
	
	userWrapper := &tokenUserWrapper{
		TokenUser: user,
		idpRole:   idpRole,
	}

	accessToken, err := jwt.GenerateAccessToken(userWrapper)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := jwt.GenerateRefreshToken(user)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

type tokenUserWrapper struct {
	jwt.TokenUser
	idpRole *string
}

func (w *tokenUserWrapper) GetIDPRole() *string {
	return w.idpRole
}

func (u *oauthUsecase) getIDPRole(userID uint64) (*string, error) {
	idpIDStr := config.Load().IdentityProviderID
	idpID, err := strconv.ParseUint(idpIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse identity provider id: %w", err)
	}

	assignedRole, err := u.assignedRolesRepository.GetByKeys(userID, idpID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
    
    return &assignedRole.RoleName, nil
}

// ─── Private helpers ──────────────────────────────────────────────────────────

// validateClientAndRedirectURI fetches the service once and validates both the
// client ID and the redirect URI, avoiding a duplicate repository call.
func (u *oauthUsecase) validateClientAndRedirectURI(clientID, redirectURI string) error {
    client, err := u.serviceRepository.GetClientByClientId(clientID)
    if err != nil || client == nil {
        return pkgErrors.ErrInvalidClientID
    }

    allowed := strings.Split(client.RedirectUrl, ",")

    for _, uri := range allowed {
        if strings.TrimSpace(uri) == redirectURI {
            return nil
        }
    }

    return pkgErrors.ErrInvalidRedirectURI
}

// validateClientCredentials checks either PKCE (code verifier) or client secret.
// It prioritizes PKCE if the session was initiated with a code challenge.
func (u *oauthUsecase) validateClientCredentials(req request.ExchangeAuthCodeRequest, session AuthSession, service *domain.Services) error {
	// PKCE Flow
	if session.CodeChallenge != nil && *session.CodeChallenge != "" {
		if req.CodeVerifier == "" {
			return fmt.Errorf("code_verifier is required for PKCE flow")
		}
		return u.verifyCodeChallenge(req.CodeVerifier, *session.CodeChallenge, session.CodeChallengeMethod)
	}

	// Client Secret Flow
	if req.ClientSecret == nil {
		return fmt.Errorf("client_secret is required")
	}
	
	return u.verifyClientSecret(service, *req.ClientSecret)
}

// verifyClientSecret compares the provided secret against the stored bcrypt hash.
func (u *oauthUsecase) verifyClientSecret(service *domain.Services, clientSecret string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(service.ClientSecret), []byte(clientSecret)); err != nil {
		return pkgErrors.ErrInvalidClientSecret
	}
	return nil
}

// verifyCodeChallenge validates the PKCE code verifier against the stored challenge.
// Supports S256 (default) and plain methods.
func (u *oauthUsecase) verifyCodeChallenge(codeVerifier, codeChallenge string, method *string) error {
	m := "S256"
	if method != nil && *method != "" {
		m = *method
	}

	switch m {
	case "S256":
		h := sha256.New()
		h.Write([]byte(codeVerifier))
		encoded := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
		if encoded != codeChallenge {
			return pkgErrors.ErrInvalidCodeVerifier
		}
	default: // "plain"
		if codeVerifier != codeChallenge {
			return pkgErrors.ErrInvalidCodeVerifier
		}
	}

	return nil
}

// replaceProfilePhoto uploads the new photo, updates the user record, and removes
// the previous photo from disk. Photo deletion is best-effort.
func (u *oauthUsecase) replaceProfilePhoto(user *domain.Users, file multipart.File, header *multipart.FileHeader) error {
	oldPhotoPath := ""
	if user.Photo != nil {
		oldPhotoPath = *user.Photo
	}

	path, err := image.ProcessImage(file, header, "public/images/profiles")
	if err != nil {
		return fmt.Errorf("failed to process profile photo: %w", err)
	}
	user.Photo = &path

	// Best-effort cleanup of the previous photo.
	if oldPhotoPath != "" {
		target := filepath.FromSlash(strings.TrimLeft(oldPhotoPath, "/\\"))
		_ = os.Remove(target) // Ignore error, as cleanup is secondary
	}

	return nil
}

// ForgotPassword handles sending username or password reset link to user's email
func (u *oauthUsecase) ForgotPassword(req request.ForgotPasswordRequest) error {
	user, err := u.oauthRepository.GetUserByEmailorUsername(req.Email)
	if err != nil {
		return nil // Avoid leaking whether user exists or not
	}

	var subject, body string
	switch req.ForgotType {
	case "username":
		subject = "Your Username Recovery"
		body, err = email.ParseTemplate(
			"internal/pkg/email/templates/forgot_username.html",
			map[string]interface{}{
				"FullName": user.FullName,
				"Username": user.Username,
			},
		)
		if err != nil {
			return err
		}

	case "password":
		token := helper.GenerateAuthCode()
		frontendURL := config.GetEnv("FRONTEND_BASE_URL", "http://localhost:3000")
		resetLink := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)

		if err := u.sessionRepository.SetSession("reset_token:"+token, user.UserId, 1*time.Hour); err != nil {
			return err
		}

		subject = "Password Reset Request"
		body, err = email.ParseTemplate(
			"internal/pkg/email/templates/forgot_password.html",
			map[string]interface{}{
				"FullName":  user.FullName,
				"ResetLink": resetLink,
			},
		)
		if err != nil {
			return err
		}
	}

	// Send email asynchronously to not block the response
	go func(emailReq request.ForgotPasswordRequest, userEmail, emailSubject, emailBody string) {
		if err := email.SendEmail(userEmail, emailSubject, emailBody); err != nil {
			log.Printf("Failed to send %s email to %s: %v", emailReq.ForgotType, userEmail, err)
		}
	}(req, user.Email, subject, body)

	return nil
}

// ResetPassword validates the reset token and updates the user's password.
func (u *oauthUsecase) ResetPassword(req request.ResetPasswordRequest) error {
	// 1. Get user ID from Redis using the token
	userIDStr, err := u.sessionRepository.GetSession("reset_token:" + req.Token)
	if err != nil {
		return pkgErrors.ErrSessionNotFound
	}
	if userIDStr == "" {
		return pkgErrors.ErrSessionExpired // or custom invalid token
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse user id from session: %w", err)
	}

	// 2. Fetch the user
	user, err := u.oauthRepository.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return pkgErrors.ErrNotFound
	}

	// 3. Hash the new password and update
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashed)

	if err := u.oauthRepository.UpdateUser(user); err != nil {
		return err
	}

	// 4. Invalidate the token so it cannot be used again
	_ = u.sessionRepository.DeleteSession("reset_token:" + req.Token)

	return nil
}
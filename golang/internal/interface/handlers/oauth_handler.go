package handler

import (
	"akastra-access/internal/app/config"
	"akastra-access/internal/app/usecases"
	"akastra-access/internal/infrastructure/middlewares"
	"akastra-access/internal/infrastructure/security/jwt"
	"akastra-access/internal/interface/http/request"
	"akastra-access/internal/interface/http/response"
	"akastra-access/internal/pkg/cookies"
	pkgErrors "akastra-access/internal/pkg/errors"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// OAuthHandler handles HTTP requests
type OAuthHandler struct {
	usecase usecases.OAuthUsecase
}

func NewOAuthHandler(u usecases.OAuthUsecase) *OAuthHandler {
	return &OAuthHandler{
		usecase: u,
	}
}

func (h *OAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var reqBody request.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		response.Forbidden(w, "invalid request body", err)
		return
	}

	if reqBody.EmailOrUsername == "" {
		response.Forbidden(w, "either email or username must be provided", nil)
		return
	}

	if reqBody.Password == "" {
		response.Forbidden(w, "password cannot be empty", nil)
		return
	}

	sessionID, user, err := h.usecase.Login(reqBody)
	if err != nil {
		response.Forbidden(w, "please check your credentials", err)
		return
	}

	cookies.SetCookie(w, "session_id", sessionID, 0)

	var photo string
	if user.Photo != nil {
		photo = *user.Photo
	}

	var phone string
	if user.Phone != nil {
		phone = *user.Phone
	}

	resp := response.LoginResponse{
		SessionID: sessionID,
		Fullname:  user.FullName,
		Email:     user.Email,
		Username:  user.Username,
		IsAdmin:   user.Admin,
		Photo:     photo, // Safely handled
		Phone:     phone, // Safely handled
	}

	response.Success(w, "success login", resp)
}

func (h *OAuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var reqBody request.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		response.Forbidden(w, "invalid request body", err)
		return
	}

	if reqBody.RefreshToken == "" {
		response.Forbidden(w, "refresh token is required", nil)
		return
	}

	accessToken, refreshToken, err := h.usecase.RefreshToken(reqBody)
	if err != nil {
		response.Forbidden(w, "invalid refresh token", err)
		return
	}

	resp := response.LoginResponseJWT{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		RefreshToken: refreshToken,
	}

	response.Success(w, "success refresh token", resp)
}

func (h *OAuthHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	codeChallenge := query.Get("code_challenge")
	var codeChallengePtr *string
	if codeChallenge != "" {
		codeChallengePtr = &codeChallenge
	}

	codeChallengeMethod := query.Get("code_challenge_method")
	var codeChallengeMethodPtr *string
	if codeChallengeMethod != "" {
		codeChallengeMethodPtr = &codeChallengeMethod
	}

	req := request.AuthorizeRequest{
		ClientID:            query.Get("client_id"),
		RedirectURI:         query.Get("redirect_uri"),
		ResponseType:        query.Get("response_type"),
		Scope:               query.Get("scope"),
		State:               query.Get("state"),
		CodeChallenge:       codeChallengePtr,
		CodeChallengeMethod: codeChallengeMethodPtr,
	}

	if err := req.Validate(); err != nil {
		response.BadRequest(w, "invalid request", err)
		return
	}

	// Try to get session_id from cookie
	var sessionID string
	cookie, err := cookies.GetCookie(r, "session_id")
	if err == nil {
		sessionID = cookie
	} else {
		// Fallback to Authorization header
		sessionID = jwt.ExtractBearer(r.Header.Get("Authorization"))
	}

	if sessionID == "" {
		h.redirectToLogin(w, r, req)
		return
	}

	authCode, err := h.usecase.Authorize(req, sessionID)
	if err != nil {
		if errors.Is(err, pkgErrors.ErrSessionNotFound) || errors.Is(err, pkgErrors.ErrSessionExpired) {
			h.redirectToLogin(w, r, req)
			return
		}

		response.Forbidden(w, "authorization failed", err)
		return
	}

	redirectURI := fmt.Sprintf("%s?code=%s&state=%s", req.RedirectURI, authCode, req.State)
	http.Redirect(w, r, redirectURI, http.StatusFound)
}

func (h *OAuthHandler) redirectToLogin(w http.ResponseWriter, r *http.Request, req request.AuthorizeRequest) {
    loginURL := config.Load().FrontendURL + "/login"

    q := url.Values{}
    q.Set("client_id", req.ClientID)
    q.Set("redirect_uri", req.RedirectURI)
    q.Set("response_type", req.ResponseType)
    q.Set("scope", req.Scope)
    q.Set("state", req.State)
	// q.Set("code_challenge", *req.CodeChallenge)
	// q.Set("code_challenge_method", *req.CodeChallengeMethod)

    fullLoginURL := fmt.Sprintf("%s?%s", loginURL, q.Encode())

    http.Redirect(w, r, fullLoginURL, http.StatusFound)
}

func (h *OAuthHandler) TokenExchange(w http.ResponseWriter, r *http.Request) {
	var reqBody request.ExchangeAuthCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		response.Forbidden(w, "invalid request body", err)
		return
	}

	if err := reqBody.Validate(); err != nil {
		response.BadRequest(w, "invalid request", err)
		return
	}

	accessToken, refreshToken, err := h.usecase.TokenExchange(reqBody)
	if err != nil {
		response.Forbidden(w, "token exchange failed", err)
		return
	}

	cookies.SetCookie(w, "refresh_token", refreshToken, 0)
	cookies.SetCookie(w, "access_token", accessToken, 0)

	resp := response.LoginResponseJWT{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		RefreshToken: refreshToken,
	}

	response.Success(w, "success exchange token", resp)
}

func (h *OAuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	user := middlewares.GetUserFromContext(r.Context())
	if user == nil {
		response.Forbidden(w, "unauthorized", nil)
		return
	}

	claims := middlewares.GetClaimsFromContext(r.Context())

	resp := response.MeResponse{
		UserID:      user.UserId,
		FullName:    user.FullName,
		Email:       user.Email,
		Username:    user.Username,
	}

	if user.Phone != nil {
		resp.Phone = user.Phone
	}
	if user.Photo != nil {
		resp.Photo = user.Photo
	}

	// Roles only exist if accessed via JWT
	if claims != nil {
		resp.ServiceName = claims.ServiceName
		resp.RoleName = claims.RoleName
	}
	
	response.Success(w, "success get user profile", resp)
}

func (h *OAuthHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	user := middlewares.GetUserFromContext(r.Context())
	if user == nil {
		response.Forbidden(w, "unauthorized", nil)
		return
	}

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.BadRequest(w, "invalid form data", err)
		return
	}

	var reqBody request.UpdateAccountRequest
	reqBody.FullName = r.FormValue("full_name")
	reqBody.Email = r.FormValue("email")
	reqBody.Username = r.FormValue("username")
	reqBody.Phone = r.FormValue("phone")
	
	password := r.FormValue("password")
	if password != "" {
		reqBody.Password = &password
	}

	passwordConfirmation := r.FormValue("password_confirmation")
	if passwordConfirmation != "" {
		reqBody.PasswordConfirmation = &passwordConfirmation
	}

	file, header, err := r.FormFile("photo")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		response.BadRequest(w, "invalid file", err)
		return
	}

	if err := reqBody.Validate(); err != nil {
		response.BadRequest(w, "invalid request", err)
		return
	}

	accessToken, refreshToken, err := h.usecase.UpdateAccount(user.UserId, reqBody, file, header)
	if err != nil {
		response.Forbidden(w, "update account failed", err)
		return
	}

	cookies.SetCookie(w, "refresh_token", refreshToken, 0)
	cookies.SetCookie(w, "access_token", accessToken, 0)

	resp := response.LoginResponseJWT{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		RefreshToken: refreshToken,
	}

	response.Success(w, "success update account", resp)
}

func (h *OAuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get session_id from cookie or header
	var sessionID string
	cookie, err := cookies.GetCookie(r, "session_id")
	if err == nil {
		sessionID = cookie
	} else {
		sessionID = jwt.ExtractBearer(r.Header.Get("Authorization"))
	}

	if sessionID != "" {
		_ = h.usecase.Logout(sessionID)
	}

	// Always clear cookie
	cookies.DeleteCookie(w, "session_id")

	response.Success(w, "success logout", nil)
}

func (h *OAuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var reqBody request.ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		response.BadRequest(w, "invalid request body", err)
		return
	}

	if err := reqBody.Validate(); err != nil {
		response.BadRequest(w, "invalid request", err)
		return
	}

	if err := h.usecase.ForgotPassword(reqBody); err != nil {
		response.Forbidden(w, "failed to process forgot password", err)
		return
	}

	response.Success(w, "if your email is registered, you will receive an email shortly", nil)
}

func (h *OAuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var reqBody request.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		response.BadRequest(w, "invalid request body", err)
		return
	}

	if err := reqBody.Validate(); err != nil {
		response.BadRequest(w, "invalid request", err)
		return
	}

	if err := h.usecase.ResetPassword(reqBody); err != nil {
		response.Forbidden(w, "invalid or expired token", err)
		return
	}

	response.Success(w, "password has been successfully updated", nil)
}

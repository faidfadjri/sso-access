package middlewares

import (
	"akastra-access/internal/infrastructure/databases/domain"
	"akastra-access/internal/infrastructure/repositories"
	"akastra-access/internal/infrastructure/security/jwt"
	"akastra-access/internal/interface/http/response"
	"akastra-access/internal/pkg/cookies"
	"context"
	"net/http"
	"strconv"
)

type contextKey string

const (
	userContextKey   contextKey = "user"
	claimsContextKey contextKey = "claims"
)

func GetUserFromContext(ctx context.Context) *domain.Users {
	user, ok := ctx.Value(userContextKey).(*domain.Users)
	if !ok {
		return nil
	}
	return user
}

// GetClaimsFromContext returns the parsed JWT claims stored by AuthJWT.
// Returns nil when the request was authenticated via session (AuthSession) instead.
func GetClaimsFromContext(ctx context.Context) *jwt.AccessClaims {
	claims, ok := ctx.Value(claimsContextKey).(*jwt.AccessClaims)
	if !ok {
		return nil
	}
	return claims
}

type Middleware struct {
	sessionRepo repositories.SessionRepository
	oauthRepo   repositories.OAuthRepository
}

func NewMiddleware(sessionRepo repositories.SessionRepository, oauthRepo repositories.OAuthRepository) *Middleware {
	return &Middleware{
		sessionRepo: sessionRepo,
		oauthRepo:   oauthRepo,
	}
}

func (m *Middleware) AuthSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Check Cookie
		sessionID, err := cookies.GetCookie(r, "session_id")
		if err != nil || sessionID == "" {
			response.Forbidden(w, "Unauthorized: session required", nil)
			return
		}

		// 2. Validate session from Redis
		userIDStr, err := m.sessionRepo.GetSession("session:" + sessionID)
		if err != nil || userIDStr == "" {
			response.Forbidden(w, "Unauthorized: invalid session", nil)
			return
		}

		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			response.Forbidden(w, "Unauthorized: invalid user id", nil)
			return
		}

		// 3. Get User
		user, err := m.oauthRepo.GetUserByID(userID)
		if err != nil || user == nil {
			response.Forbidden(w, "Unauthorized: user not found", nil)
			return
		}

		// 4. Set user to context
		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) AuthJWT(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var tokenString string
        var err error

        // 1. Auth Header Checking
        authHeader := r.Header.Get("Authorization")
        if authHeader != "" {
            tokenString = jwt.ExtractBearer(authHeader)
        }

        // 2. Cookie Checking
        if tokenString == "" {
            tokenString, err = cookies.GetCookie(r, "access_token")
            if err != nil || tokenString == "" {
                response.Forbidden(w, "Unauthorized: token required", nil)
                return
            }
        }

        // 3. Validate Token
        claims, err := jwt.ValidateAccessToken(tokenString)
        if err != nil {
            response.Forbidden(w, "Unauthorized: invalid access token", nil)
            return
        }

        // 4. Get User
        user, err := m.oauthRepo.GetUserByID(claims.UserID)
        if err != nil || user == nil {
            response.Forbidden(w, "Unauthorized: user not found", nil)
            return
        }

        // 5. Store both the DB user and the JWT claims in context.
        // Claims carry service_name and role_name that the plain Users record lacks.
        ctx := context.WithValue(r.Context(), userContextKey, user)
        ctx = context.WithValue(ctx, claimsContextKey, claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// AuthIDP fallback checks JWT first, then falls back to session check
func (m *Middleware) AuthIDP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Try JWT
		var tokenString string
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			tokenString = jwt.ExtractBearer(authHeader)
		}
		if tokenString == "" {
			tokenString, _ = cookies.GetCookie(r, "access_token")
		}

		if tokenString != "" {
			claims, err := jwt.ValidateAccessToken(tokenString)
			if err == nil {
				user, err := m.oauthRepo.GetUserByID(claims.UserID)
				if err == nil && user != nil {
					ctx := context.WithValue(r.Context(), userContextKey, user)
					ctx = context.WithValue(ctx, claimsContextKey, claims)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}
		}

		// 2. Fallback to Session
		sessionID, err := cookies.GetCookie(r, "session_id")
		if err == nil && sessionID != "" {
			userIDStr, err := m.sessionRepo.GetSession("session:" + sessionID)
			if err == nil && userIDStr != "" {
				userID, err := strconv.ParseUint(userIDStr, 10, 64)
				if err == nil {
					user, err := m.oauthRepo.GetUserByID(userID)
					if err == nil && user != nil {
						ctx := context.WithValue(r.Context(), userContextKey, user)
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
				}
			}
		}

		// 3. Unauthorized if both failed
		response.Forbidden(w, "Unauthorized: valid session or token required", nil)
	})
}


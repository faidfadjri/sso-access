package middlewares

import (
	"akastra-access/internal/infrastructure/security/jwt"
	"akastra-access/internal/interface/http/response"
	"akastra-access/internal/pkg/cookies"
	"context"
	"net/http"
)

// RequireIDPRole creates a middleware that checks if the user has at least one of the required IDP roles.
func (m *Middleware) RequireIDPRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaimsFromContext(r.Context())
			
			// If no claims in context, try to extract from cookie/header (like AuthJWT)
			if claims == nil {
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
				}

				// 3. Validate Token if found
				if tokenString != "" && err == nil {
					validClaims, err := jwt.ValidateAccessToken(tokenString)
					if err == nil {
						claims = validClaims
						// Optional: Update context for downstream handlers?
						ctx := context.WithValue(r.Context(), claimsContextKey, claims)
						r = r.WithContext(ctx)
					}
				}
			}

			if claims == nil {
				response.Forbidden(w, "Unauthorized: no valid token found", nil)
				return
			}

			if claims.IDPRole == nil {
				response.Forbidden(w, "Forbidden: no IDP role", nil)
				return
			}

			userRole := *claims.IDPRole
			for _, role := range roles {
				if userRole == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			response.Forbidden(w, "Forbidden: insufficient permissions", nil)
		})
	}
}

// Admin middleware allows access to users with "Admin" or "Super Admin" IDP role.
func (m *Middleware) AdminMiddleware(next http.Handler) http.Handler {
	return m.RequireIDPRole("Admin", "Super Admin")(next)
}

// SuperAdmin middleware allows access to users with "Super Admin" IDP role only.
func (m *Middleware) SuperAdminMiddleware(next http.Handler) http.Handler {
	return m.RequireIDPRole("Super Admin")(next)
}

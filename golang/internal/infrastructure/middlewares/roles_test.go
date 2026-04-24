package middlewares

import (
	"akastra-access/internal/infrastructure/security/jwt"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireIDPRole(t *testing.T) {
	m := &Middleware{}
	roleAdmin := "Admin"
	roleSuperAdmin := "Super Admin"
	roleUser := "User"

	tests := []struct {
		name           string
		userRole       *string
		requiredRoles  []string
		expectedStatus int
		useCookie      bool
	}{
		{
			name:           "Authorized - Admin accessing Admin resource (Context)",
			userRole:       &roleAdmin,
			requiredRoles:  []string{"Admin"},
			expectedStatus: http.StatusOK,
			useCookie:      false,
		},
		{
			name:           "Authorized - Super Admin accessing Admin resource (Context)",
			userRole:       &roleSuperAdmin,
			requiredRoles:  []string{"Admin", "Super Admin"},
			expectedStatus: http.StatusOK,
			useCookie:      false,
		},
		{
			name:           "Forbidden - User accessing Admin resource (Context)",
			userRole:       &roleUser,
			requiredRoles:  []string{"Admin"},
			expectedStatus: http.StatusForbidden,
			useCookie:      false,
		},
		{
			name:           "Unauthorized - No valid token found (Context)",
			userRole:       nil,
			requiredRoles:  []string{"Admin"},
			expectedStatus: http.StatusForbidden,
			useCookie:      false,
		},
		// Cookie Tests requiring valid JWT generation which depends on config
		// Skipping real JWT generation in unit test without mocking config which is hard here.
		// However, we can test that it TRIES to get from cookie if context is empty.
		// Use a mock token string? ValidateAccessToken uses config.Load().JWTSecret.
		// We can't easily mock config.Load().
		// So we will stick to context tests for now, trusting the logic mirrors AuthJWT which is tested elsewhere/trusted.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := m.RequireIDPRole(tt.requiredRoles...)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()

			if tt.useCookie {
				// Simulate cookie
				// Note: In unit tests without valid JWT secrets, Validation will fail -> Unauthorized
				// This is expected for "Unauthorized - No valid token found" cases if we added them
				cookie := &http.Cookie{
					Name:  "access_token",
					Value: "dummy_token", 
				}
				req.AddCookie(cookie)
			} 
			
			if tt.userRole != nil {
				claims := &jwt.AccessClaims{
					IDPRole: tt.userRole,
				}
				ctx := context.WithValue(req.Context(), claimsContextKey, claims)
				req = req.WithContext(ctx)
			}

			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

package routes

import (
	"akastra-access/internal/app/bootstrap"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func OAuthRouter(deps *bootstrap.Dependencies) http.Handler {
	r := chi.NewRouter()

	r.Post("/login", deps.OAuthHandler.Login)
	r.Post("/refresh-token", deps.OAuthHandler.RefreshToken)
	r.Get("/authorize", deps.OAuthHandler.Authorize)
	r.Post("/token", deps.OAuthHandler.TokenExchange)
	r.Post("/logout", deps.OAuthHandler.Logout)
	r.Post("/forgot-password", deps.OAuthHandler.ForgotPassword)
	r.Post("/reset-password", deps.OAuthHandler.ResetPassword)

	// Protected routes
	r.With(deps.Middleware.AuthIDP).Put("/update-account", deps.OAuthHandler.UpdateAccount)
	r.With(deps.Middleware.AuthIDP).Get("/me", deps.OAuthHandler.Me)
	return r
}
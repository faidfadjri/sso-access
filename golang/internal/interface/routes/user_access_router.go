package routes

import (
	"akastra-access/internal/app/bootstrap"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func UserAccessRouter(deps *bootstrap.Dependencies) http.Handler {
	r := chi.NewRouter()

	r.Use(deps.Middleware.AuthSession)

	r.Get("/", deps.UserAccessHandler.GetUserAccesses)
	r.Post("/", deps.UserAccessHandler.CreateUserAccess)
	r.Get("/{id}", deps.UserAccessHandler.GetUserAccessByID)
	r.Put("/{id}", deps.UserAccessHandler.UpdateUserAccess)
	r.Delete("/{id}", deps.UserAccessHandler.DeleteUserAccess)

	return r
}

package routes

import (
	"akastra-access/internal/app/bootstrap"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func UserRouter(deps *bootstrap.Dependencies) http.Handler {
	r := chi.NewRouter()

	r.Use(deps.Middleware.AuthSession)
	// r.Use(deps.Middleware.Admin)	

	r.Get("/", deps.UserHandler.GetUsers)
	r.Post("/", deps.UserHandler.CreateUser)
	r.Put("/{id}", deps.UserHandler.UpdateUser)
	r.Delete("/{id}", deps.UserHandler.DeleteUser)
	r.Delete("/batch", deps.UserHandler.BatchDeleteUser)

	return r
}

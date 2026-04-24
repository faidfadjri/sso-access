package routes

import (
	"akastra-access/internal/app/bootstrap"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func AssignedRolesRouter(deps *bootstrap.Dependencies) http.Handler {
	r := chi.NewRouter()

	r.Use(deps.Middleware.AuthSession)

	r.Get("/", deps.AssignedRolesHandler.GetList)
	r.Post("/", deps.AssignedRolesHandler.AssignRole)
	r.Delete("/", deps.AssignedRolesHandler.RemoveRole)

	return r
}

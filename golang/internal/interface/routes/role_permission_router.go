package routes

import (
	"akastra-access/internal/app/bootstrap"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RolePermissionRouter(deps *bootstrap.Dependencies) http.Handler {
	r := chi.NewRouter()

	r.Use(deps.Middleware.AuthSession)
	// r.Use(deps.Middleware.Admin)

	r.Get("/", deps.RolePermissionHandler.GetList)
	r.Post("/", deps.RolePermissionHandler.AssignPermission)
	r.Delete("/", deps.RolePermissionHandler.RemovePermission)

	return r
}

package routes

import (
	"akastra-access/internal/app/bootstrap"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func PermissionRouter(deps *bootstrap.Dependencies) http.Handler {
	r := chi.NewRouter()

	r.Use(deps.Middleware.AuthSession)
	// r.Use(deps.Middleware.Admin)

	r.Get("/", deps.PermissionHandler.GetPermissions)
	r.Post("/", deps.PermissionHandler.CreatePermission)
	r.Get("/{id}", deps.PermissionHandler.GetPermissionByID)
	r.Put("/{id}", deps.PermissionHandler.UpdatePermission)
	r.Delete("/{id}", deps.PermissionHandler.DeletePermission)

	return r
}

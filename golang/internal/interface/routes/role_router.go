package routes

import (
	"akastra-access/internal/app/bootstrap"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RoleRouter(deps *bootstrap.Dependencies) http.Handler {
	r := chi.NewRouter()

	r.Use(deps.Middleware.AuthSession)
	// r.Use(deps.Middleware.Admin)

	r.Get("/", deps.RoleHandler.GetRoles)
	r.Post("/", deps.RoleHandler.CreateRole)
	r.Get("/{id}", deps.RoleHandler.GetRoleByID)
	r.Put("/{id}", deps.RoleHandler.UpdateRole)
	r.Delete("/{id}", deps.RoleHandler.DeleteRole)


	r.Route("/assign", func(r chi.Router) {
		r.Get("/", deps.AssignedRolesHandler.GetList)
		r.Post("/", deps.AssignedRolesHandler.AssignRole)
		r.Delete("/", deps.AssignedRolesHandler.RemoveRole)
	})

	return r
}

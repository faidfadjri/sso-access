package routes

import (
	"akastra-access/internal/app/bootstrap"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func ServiceRouter(deps *bootstrap.Dependencies) http.Handler {
	r := chi.NewRouter()

	r.Post("/clients", deps.ServiceHandler.CreateClient)
	r.Get("/clients", deps.ServiceHandler.GetClients)
	r.Delete("/clients/{id}", deps.ServiceHandler.DeleteClientById)
	r.Put("/clients/{id}", deps.ServiceHandler.UpdateClient)
	return r
}
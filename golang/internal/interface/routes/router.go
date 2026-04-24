package routes

import (
	"akastra-access/internal/app/bootstrap"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func InitRouter(deps *bootstrap.Dependencies) http.Handler {
	r := chi.NewRouter()

	// 1. CORS Configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"https://access.akastra.id",
			"https://www.access.akastra.id",
			"http://localhost:3000",
		},
		AllowedMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS",
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
		},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// 2. Health Check
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from Akastra Access Provider!"))
	})

	// 3. Static Files
	fileServer := http.FileServer(http.Dir("./public"))
	r.Handle("/public/*", http.StripPrefix("/public", fileServer))

	// 4. API Routes (v1)
	const prefix = "/api/v1"
	r.Route(prefix, func(api chi.Router) {
		
		// Public Routes
		api.Mount("/oauth", OAuthRouter(deps))

		// Protected Routes (Auth + Admin check applied internally or per handler)
		api.Mount("/users", UserRouter(deps))
		api.Mount("/roles", RoleRouter(deps))
		api.Mount("/permissions", PermissionRouter(deps))
		
		// Access Management
		api.Mount("/users/access", UserAccessRouter(deps))
		api.Mount("/roles/permissions", RolePermissionRouter(deps))

		// Complex Service Routes
		api.Route("/service", func(router chi.Router) {
			router.Use(deps.Middleware.AuthSession)
			router.Use(deps.Middleware.SuperAdminMiddleware)
			router.Mount("/", ServiceRouter(deps))
		})
	})


	return r
}

package server

import (
	"net/http"

	"github.com/FactoryCICD/factory-api/internal/datastore"
	"github.com/FactoryCICD/factory-api/internal/modules/projects"
	dbcontext "github.com/FactoryCICD/factory-api/pkg/db"
	"github.com/FactoryCICD/factory-api/pkg/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

// Build creates a new router and adds the routes to it
func Build(logger log.Logger, db *dbcontext.DB) *chi.Mux {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Routes
	r.Get("/", rootRoute)
	r.Get("/health", healthCheck)

	// Handlers
	r.Route("/v1", func(r chi.Router) {
		r.Mount("/projects", projects.Routes(projects.NewController(datastore.NewDatastore(db), logger)))
	})

	// Print routes
	printEstablishedRoutes(r, logger)
	return r
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	// TODO: Return more information about the health of the service
	message := struct {
		Message string `json:"message"`
	}{Message: "OK"}
	render.JSON(w, r, message)
}

func rootRoute(w http.ResponseWriter, r *http.Request) {
	message := struct {
		Message string `json:"message"`
	}{Message: "Hello World"}
	render.JSON(w, r, message)
}

func printEstablishedRoutes(r *chi.Mux, logger log.Logger) {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		logger.Info("Route Established: ", method, " - ", route)
		return nil
	}
	if err := chi.Walk(r, walkFunc); err != nil {
		logger.Error("Error while printing routes: ", err.Error())
	}
}

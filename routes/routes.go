package routes

import (
	"net/http"
	"portfolio-api/controllers"
	"portfolio-api/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Routes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Project Rotues
	router.Get("/api/v1/projects", controllers.GetProjects)
	router.Get("/api/v1/projects/project/{id}", controllers.GetProjectById)
	router.Route("/api/v1/projects/project", func(router chi.Router) {
		router.Use(middlewares.IsAuthorized)
		router.Post("/", controllers.CreateProject)
		router.Put("/{id}", controllers.UpdateProject)
		router.Delete("/{id}", controllers.DeleteProject)
	})

	// Auth Routes
	// router.Post("/api/v1/auth/signup", controllers.Singup)
	router.Post("/api/v1/auth/login", controllers.Login)
	router.Route("/api/v1/auth/user", func(router chi.Router) {
		router.Get("/bytoken", controllers.GetUserByToken)
	})

	return router
}

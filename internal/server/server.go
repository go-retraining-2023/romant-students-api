package server

import (
	"net/http"

	handlers "github.com/RomanTykhyi/students-api/internal/server/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func StartServer() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/api/v1/students", func(r chi.Router) {
		r.Get("/", handlers.QueryStudents)
		r.Post("/", handlers.CreateStudent)

		// subroute
		r.Route("/{studentId}", func(r chi.Router) {
			r.Get("/", handlers.GetStudent)
			r.Put("/", handlers.UpdateStudent)
			r.Delete("/", handlers.DeleteStudent)
		})
	})

	http.ListenAndServe(":8081", r)
}

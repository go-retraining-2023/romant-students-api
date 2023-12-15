package server

import (
	"net/http"

	handlers "github.com/RomanTykhyi/students-api/internal/server/handlers"
	utils "github.com/RomanTykhyi/students-api/internal/server/utils"
)

func StartServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/students", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.QueryStudents(w, r)
		case "POST":
			handlers.CreateStudent(w, r)
		default:
			utils.WriteString(w, "Unsupported method")
		}
	})

	mux.HandleFunc("/api/students/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetStudent(w, r)
		case http.MethodPut:
			handlers.UpdateStudent(w, r)
		case http.MethodDelete:
			handlers.DeleteStudent(w, r)
		default:
			utils.WriteString(w, "Unsupported method")
		}
	})

	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	server.ListenAndServe()
}

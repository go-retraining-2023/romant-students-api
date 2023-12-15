package server

import (
	"encoding/json"
	"net/http"

	"github.com/RomanTykhyi/students-api/internal/data"
	"github.com/RomanTykhyi/students-api/internal/di"
	"github.com/RomanTykhyi/students-api/internal/models"
)

func QueryStudents(w http.ResponseWriter, r *http.Request) []models.Student {
	repo, err := di.GetAppContainer().Resolve("students-store")
	if err != nil {
		panic("Cannot get students repository")
	}
	studentsRepo := repo.(data.StudentsStore)

	students := studentsRepo.QueryStudents()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)

	return students
}

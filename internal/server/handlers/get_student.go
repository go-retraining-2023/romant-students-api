package server

import (
	"encoding/json"
	"net/http"

	"github.com/RomanTykhyi/students-api/internal/data"
	"github.com/RomanTykhyi/students-api/internal/di"
	utils "github.com/RomanTykhyi/students-api/internal/server/utils"
	"github.com/go-chi/chi/v5"
)

func GetStudent(w http.ResponseWriter, r *http.Request) {
	repo, err := di.GetAppContainer().Resolve("students-store")
	if err != nil {
		panic("Cannot get students repository")
	}

	studentsRepo := repo.(data.StudentsStore)

	studentId := chi.URLParam(r, "studentId")

	student := studentsRepo.GetStudent(studentId)
	if student == nil {
		w.WriteHeader(http.StatusNotFound)
		utils.WriteString(w, "Student not found.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

package server

import (
	"net/http"

	"github.com/RomanTykhyi/students-api/internal/data"
	"github.com/RomanTykhyi/students-api/internal/di"
	utils "github.com/RomanTykhyi/students-api/internal/server/utils"
	"github.com/go-chi/chi/v5"
)

func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	repo, err := di.GetAppContainer().Resolve("students-store")
	if err != nil {
		panic("Cannot get students repository")
	}
	studentsRepo := repo.(data.StudentsStore)

	studentId := chi.URLParam(r, "studentId")

	deleted, err := studentsRepo.DeleteStudent(studentId)
	if err != nil {
		utils.WriteMessageResponse(w, "Error while deleting.", http.StatusInternalServerError)
		return
	}

	if deleted {
		utils.WriteMessageResponse(w, "Deleted successfully.", http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

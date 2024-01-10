package server

import (
	"net/http"

	"github.com/RomanTykhyi/students-api/internal/data"
	"github.com/RomanTykhyi/students-api/internal/di"
	utils "github.com/RomanTykhyi/students-api/internal/server/utils"
)

func QueryStudents(w http.ResponseWriter, r *http.Request) {
	repo, err := di.GetAppContainer().Resolve("students-store")
	if err != nil {
		panic("Cannot get students repository")
	}
	studentsRepo := repo.(data.StudentsStore)

	students := studentsRepo.QueryStudents()
	if students == nil {
		utils.WriteMessageResponse(w, "Error occured", http.StatusInternalServerError)
		return
	}

	utils.WriteJsonResponse(w, students)
}

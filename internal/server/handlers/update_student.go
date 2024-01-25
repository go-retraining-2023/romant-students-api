package server

import (
	"net/http"

	"github.com/RomanTykhyi/students-api/internal/data"
	"github.com/RomanTykhyi/students-api/internal/di"
	"github.com/RomanTykhyi/students-api/internal/models"
	utils "github.com/RomanTykhyi/students-api/internal/server/utils"
	"github.com/go-chi/chi/v5"
)

func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	repo, err := di.GetAppContainer().Resolve("students-store")
	if err != nil {
		panic("Cannot get students repository")
	}
	studentsRepo := repo.(data.StudentsStore)

	if err := r.ParseForm(); err != nil {
		utils.WriteMessageResponse(w, "Failed to parse form", http.StatusInternalServerError)
		return
	}

	fullName := r.PostForm["FullName"][0]

	studentId := chi.URLParam(r, "studentId")

	student := models.Student{
		PartitionId: "students",
		Id:          studentId,
		FullName:    fullName,
	}

	_, err = studentsRepo.UpdateStudent(&student)
	if err != nil {
		utils.WriteMessageResponse(w, "Failed to update student.", http.StatusInternalServerError)
		return
	}

	utils.WriteJsonResponse(w, student)
}

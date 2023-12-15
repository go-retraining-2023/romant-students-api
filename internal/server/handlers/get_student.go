package server

import (
	"encoding/json"
	"net/http"

	"github.com/RomanTykhyi/students-api/internal/data"
	"github.com/RomanTykhyi/students-api/internal/di"
	"github.com/RomanTykhyi/students-api/internal/models"
	utils "github.com/RomanTykhyi/students-api/internal/server/utils"
)

func GetStudent(w http.ResponseWriter, r *http.Request) *models.Student {
	repo, err := di.GetAppContainer().Resolve("students-store")
	if err != nil {
		panic("Cannot get students repository")
	}
	studentsRepo := repo.(data.StudentsStore)

	studentId, err := utils.RetrieveStudentId(w, r)
	if err != nil {
		return nil
	}

	student := studentsRepo.GetStudent(studentId.String())
	if student == nil {
		w.WriteHeader(http.StatusNotFound)
		utils.WriteString(w, "Student not found.")
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)

	return student
}

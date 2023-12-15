package server

import (
	"net/http"

	"github.com/RomanTykhyi/students-api/internal/data"
	"github.com/RomanTykhyi/students-api/internal/di"
	utils "github.com/RomanTykhyi/students-api/internal/server/utils"
)

func DeleteStudent(w http.ResponseWriter, r *http.Request) bool {
	repo, err := di.GetAppContainer().Resolve("students-store")
	if err != nil {
		panic("Cannot get students repository")
	}
	studentsRepo := repo.(data.StudentsStore)

	studentId, err := utils.RetrieveStudentId(w, r)
	if err != nil {
		return false
	}

	return studentsRepo.DeleteStudent(studentId.String())
}

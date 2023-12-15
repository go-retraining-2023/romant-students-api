package server

import (
	"encoding/json"
	"net/http"

	"github.com/RomanTykhyi/students-api/internal/data"
	"github.com/RomanTykhyi/students-api/internal/di"
	"github.com/RomanTykhyi/students-api/internal/models"
	utils "github.com/RomanTykhyi/students-api/internal/server/utils"
)

func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	repo, err := di.GetAppContainer().Resolve("students-store")
	if err != nil {
		panic("Cannot get students repository")
	}
	studentsRepo := repo.(data.StudentsStore)

	if err := r.ParseForm(); err != nil {
		w.Write([]byte("Failed to parse from..."))
	}

	fullName := r.PostForm["FullName"][0]

	studentId, err := utils.RetrieveStudentId(w, r)
	if err != nil {
		return
	}

	student := models.Student{
		PartitionId: "students",
		Id:          studentId.String(),
		FullName:    fullName,
	}

	studentsRepo.UpdateStudent(&student)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

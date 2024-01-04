package server

import (
	"encoding/csv"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/RomanTykhyi/students-api/internal/data"
	"github.com/RomanTykhyi/students-api/internal/di"
	utils "github.com/RomanTykhyi/students-api/internal/server/utils"
)

func ExportStudents(w http.ResponseWriter, r *http.Request) {
	repo, err := di.GetAppContainer().Resolve("students-store")
	if err != nil {
		panic("Cannot get students repository")
	}
	studentsRepo := repo.(data.StudentsStore)

	// get students
	students := studentsRepo.QueryStudents()

	// parse to csv
	records := [][]string{{"Id", "FullName"}}

	for _, student := range students {
		studentRecord := []string{student.Id, student.FullName}
		records = append(records, studentRecord)
	}

	// trying to create temporary file
	f, err := os.CreateTemp("", "students_*.csv")
	if err != nil {
		utils.WriteError(w, "Error creating file.", http.StatusInternalServerError)
	}

	defer f.Close()

	// write our records
	csvWriter := csv.NewWriter(f)
	csvWriter.WriteAll(records)

	f.Seek(0, io.SeekStart)

	// downloading the file
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(f.Name()))
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

	io.Copy(w, f)
	csvWriter.Flush()
}

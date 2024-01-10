package server

import (
	"encoding/csv"
	"net/http"

	"github.com/RomanTykhyi/students-api/internal/data"
	"github.com/RomanTykhyi/students-api/internal/models"
	utils "github.com/RomanTykhyi/students-api/internal/server/utils"

	"github.com/RomanTykhyi/students-api/internal/di"
)

const maxUploadSize = 10 * 1024 * 1024 // 10 mb

func ImportStudents(w http.ResponseWriter, r *http.Request) {
	// getting repo from di container
	repo, err := di.GetAppContainer().Resolve("students-store")
	if err != nil {
		utils.WriteMessageResponse(w, "Could not get the students store.", http.StatusInternalServerError)
		return
	}
	studentsRepo := repo.(data.StudentsStore)

	// parsing multipart form file
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		utils.WriteMessageResponse(w, "Could not parse multipart form", http.StatusInternalServerError)
		return
	}

	// parse and validate file and post parameters
	file, fileHeader, err := r.FormFile("uploadFile")
	if err != nil {
		utils.WriteMessageResponse(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// validating the file size
	if fileHeader.Size > maxUploadSize {
		utils.WriteMessageResponse(w, "The file is too big. Max size is 10 mb", http.StatusBadRequest)
		return
	}

	// checking the content-type
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType != "text/csv" {
		utils.WriteMessageResponse(w, "Invalid content-type. Only csv import allowed", http.StatusBadRequest)
		return
	}

	// read csv records
	fileReader := csv.NewReader(file)
	records, error := fileReader.ReadAll()
	if error != nil {
		utils.WriteMessageResponse(w, "Error while reading file.", http.StatusInternalServerError)
		return
	}
	records = records[1:]

	// create student and save
	for _, rec := range records {
		student := models.Student{
			Id:       rec[0],
			FullName: rec[1],
		}

		err := studentsRepo.PutStudent(&student)
		if err != nil {
			utils.WriteMessageResponse(w, "Error while inserting student", http.StatusInternalServerError)
		}
	}
}

package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"log"
	"mime"
	"mime/multipart"
	"net/http"

	"github.com/RomanTykhyi/students-api/internal/common"
	"github.com/RomanTykhyi/students-api/internal/data"
	"github.com/RomanTykhyi/students-api/internal/di"
	"github.com/RomanTykhyi/students-api/internal/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	repo, err := di.GetAppContainer().Resolve("students-store")
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Can't resolve data store.",
		}, nil
	}
	studentsRepo := repo.(data.StudentsStore)

	contentType := request.Headers["Content-Type"]

	// creating the multipart reader
	mr := multipart.NewReader(bytes.NewReader([]byte(request.Body)), extractBoundary(contentType))
	log.Println("Multipart reader created.")

	// get the part
	part, err := mr.NextPart()
	if err == http.ErrMissingFile {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Missing file.",
		}, nil
	} else if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error reading the form data.",
		}, nil
	}

	// determine that we actually have the file sent
	if part.FileName() == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error reading the file.",
		}, nil
	}
	log.Printf("Got the %v file.", part.FileName())

	reader := csv.NewReader(part)
	log.Println("Reading the csv data.")

	records, err := reader.ReadAll()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error parsing the csv data.",
		}, nil
	}

	records = records[1:]

	// create student and save
	for _, rec := range records {
		student := models.Student{
			Id:       rec[0],
			FullName: rec[1],
		}

		err := studentsRepo.PutStudent(&student)
		log.Printf("Saving student:%v", student.FullName)

		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       "Error saving the student.",
			}, nil
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Import successful.",
	}, nil
}

func init() {
	common.Init()
}

func main() {
	lambda.Start(HandleRequest)
}

func extractBoundary(contentType string) string {
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return ""
	}
	return params["boundary"]
}

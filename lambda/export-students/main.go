package main

import (
	"context"
	"encoding/csv"
	"net/http"
	"strings"

	"github.com/RomanTykhyi/students-api/internal/common"
	"github.com/RomanTykhyi/students-api/internal/data"
	"github.com/RomanTykhyi/students-api/internal/di"
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

	// get students
	students, err := studentsRepo.QueryStudents()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error querying students.",
		}, nil
	}

	// convert to csv
	records := [][]string{{"Id", "FullName"}}

	for _, student := range students {
		studentRecord := []string{student.Id, student.FullName}
		records = append(records, studentRecord)
	}

	// write our records
	var csvContent strings.Builder
	csvWriter := csv.NewWriter(&csvContent)
	err = csvWriter.WriteAll(records)
	csvWriter.Flush()

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error on writing the csv data.",
		}, nil
	}

	// set response headers
	headers := map[string]string{
		"Content-Disposition": "attachment; filename=myfile.csv",
		"Content-Type":        "application/octet-stream",
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       csvContent.String(),
	}, nil
}

func init() {
	common.Init()
}

func main() {
	lambda.Start(HandleRequest)
}

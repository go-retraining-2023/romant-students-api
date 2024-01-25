package main

import (
	"context"
	"encoding/json"
	"net/http"

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

	students, err := studentsRepo.QueryStudents()
	if err != nil {

	}

	if students == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error querying students.",
		}, nil
	}

	studentsBytes, err := json.Marshal(students)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error marshaling students.",
		}, nil
	}

	studentsJson := string(studentsBytes)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       studentsJson,
	}, nil
}

func init() {
	common.Init()
}

func main() {
	lambda.Start(HandleRequest)
}

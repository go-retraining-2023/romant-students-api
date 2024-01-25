package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/RomanTykhyi/students-api/internal/common"
	"github.com/RomanTykhyi/students-api/internal/data"
	"github.com/RomanTykhyi/students-api/internal/di"
	"github.com/RomanTykhyi/students-api/internal/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
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

	formParams, err := url.ParseQuery(request.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid request body",
		}, nil
	}

	fullName := formParams.Get("FullName")

	student := models.Student{
		PartitionId: "students",
		Id:          uuid.New().String(),
		FullName:    fullName,
	}

	err = studentsRepo.PutStudent(&student)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error creating student.",
		}, nil
	}

	studentsBytes, err := json.Marshal(student)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error marshaling student.",
		}, nil
	}

	studentJson := string(studentsBytes)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       studentJson,
	}, nil
}

func init() {
	common.Init()
}

func main() {
	lambda.Start(HandleRequest)
}

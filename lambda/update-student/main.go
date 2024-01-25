package main

import (
	"context"
	"encoding/json"
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

	if request.Body == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Request body is empty.",
		}, nil
	}

	// get FullName from body
	var student models.Student
	err = json.Unmarshal([]byte(request.Body), &student)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Error unmarshaling student.",
		}, nil
	}

	// get student id
	studentIdRaw, found := request.PathParameters["studentId"]
	if !found {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Can't retrieve parameter from path.",
		}, nil
	}

	student.Id = studentIdRaw
	student.PartitionId = "students"
	_, err = studentsRepo.UpdateStudent(&student)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error updating student.",
		}, nil
	}

	studentsBytes, err := json.Marshal(student)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error marshaling student.",
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

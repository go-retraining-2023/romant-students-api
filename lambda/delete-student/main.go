package main

import (
	"context"
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

	studentIdRaw, found := request.PathParameters["studentId"]
	if !found {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Can't retrieve parameter from path.",
		}, nil
	}

	deleted, err := studentsRepo.DeleteStudent(studentIdRaw)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while deleting.",
		}, nil
	}

	if deleted {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNoContent,
			Body:       "Deleted successfully.",
		}, nil
	} else {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Looks like there is no such student to delete.",
		}, nil
	}
}

func init() {
	common.Init()
}

func main() {
	lambda.Start(HandleRequest)
}

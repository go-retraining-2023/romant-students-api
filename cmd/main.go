package main

import (
	"github.com/RomanTykhyi/students-api/internal/data"
	"github.com/RomanTykhyi/students-api/internal/di"
	"github.com/RomanTykhyi/students-api/internal/server"
)

func main() {

	appContainer := di.GetAppContainer()

	dynamoClient := data.CreateLocalClient()
	appContainer.Register("dynamo-client", dynamoClient)                // dynamodb client
	appContainer.Register("students-store", &data.StudentsRepository{}) // repository

	server.StartServer(8081)
}

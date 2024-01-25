package common

import (
	"github.com/RomanTykhyi/students-api/config"
	"github.com/RomanTykhyi/students-api/internal/data"
	"github.com/RomanTykhyi/students-api/internal/di"
)

func Init() {
	appConfig := config.CreateAppConfig()
	awsConfig := config.CreateAWSConfig()
	container := di.GetAppContainer()
	ConfigureDIContainer(appConfig.DynamoDbUrl, awsConfig, container)
}

func ConfigureDIContainer(dynamoDbUrl string, awsConfig *config.AwsConfig, appContainer *di.Container) {
	dynamoClient, err := data.CreateLocalClient(dynamoDbUrl, awsConfig)
	if err != nil {
		panic("Error creating dynamodb client.")
	}

	studentsRepo := data.NewStudentsRepository(dynamoClient)

	appContainer.Register("dynamo-client", dynamoClient)  // dynamodb client
	appContainer.Register("students-store", studentsRepo) // repository
}

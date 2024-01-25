package config

import (
	"os"
	"strconv"
)

var appConfig *ApplicationConfig

type ApplicationConfig struct {
	HttpPort    int
	DynamoDbUrl string
}

func CreateAppConfig() *ApplicationConfig {
	appConfig = &ApplicationConfig{
		HttpPort:    convertStringToInteger("HTTP_PORT", 8081),
		DynamoDbUrl: getValueFromEnvOrDefault("DYNAMODB_URL", "http://localhost:8127"),
	}

	return appConfig
}

func GetAppConfig() *ApplicationConfig {
	if appConfig == nil {
		CreateAppConfig()
	}

	return appConfig
}

func convertStringToInteger(value string, defaultValue int) int {
	result, err := strconv.Atoi(value)

	if err != nil {
		return defaultValue
	}

	return result
}

type AwsConfig struct {
	AwsDefaultRegion   string
	AwsAccessKey       string
	AwsSecretAccessKey string
	AwsSessionToken    string
}

func CreateAWSConfig() *AwsConfig {
	return &AwsConfig{
		AwsDefaultRegion:   getValueFromEnvOrDefault("AWS_DEFAULT_REGION", "localhost"),
		AwsAccessKey:       getValueFromEnvOrDefault("AWS_ACCESS_KEY_ID", "abcd"),
		AwsSecretAccessKey: getValueFromEnvOrDefault("AWS_SECRET_ACCESS_KEY", "a1b2c3"),
		AwsSessionToken:    getValueFromEnvOrDefault("AWS_SESSION_TOKEN", ""),
	}
}

func getValueFromEnvOrDefault(envKey string, defaultValue string) string {
	envVarValue := os.Getenv(envKey)
	if envVarValue == "" {
		return defaultValue
	}

	return envVarValue
}

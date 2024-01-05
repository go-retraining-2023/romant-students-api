package data

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func CreateLocalClient() *dynamodb.Client {

	dbUrl := os.Getenv("DYNAMODB_URL")
	if dbUrl == "" {
		dbUrl = "http://localhost:8127"
	}

	awsRegion := os.Getenv("AWS_DEFAULT_REGION")
	if awsRegion == "" {
		awsRegion = "localhost"
	}

	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	if awsAccessKey == "" {
		awsAccessKey = "abcd"
	}

	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if awsSecretAccessKey == "" {
		awsSecretAccessKey = "a1b2c3"
	}

	awsSessionToken := os.Getenv("AWS_SESSION_TOKEN")

	log.Printf("Dynamo url:%v", dbUrl)
	endpointResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: dbUrl,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithEndpointResolverWithOptions(endpointResolver),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     awsAccessKey,
				SecretAccessKey: awsSecretAccessKey,
				SessionToken:    awsSessionToken,
				Source:          "Mock credentials used above for local instance",
			},
		}))

	fmt.Printf("Config is: %v", cfg)

	if err != nil {
		log.Fatal(err)
	}

	return dynamodb.NewFromConfig(cfg)
}

func CreateClient() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile("roman-aws"),
	)

	if err != nil {
		log.Fatal(err)
	}

	return dynamodb.NewFromConfig(cfg)
}

func GetTables(dynamoDbClient *dynamodb.Client) []string {
	tablesOutput, err := dynamoDbClient.ListTables(
		context.TODO(),
		&dynamodb.ListTablesInput{})

	if err != nil {
		log.Fatal(err)
	}

	return tablesOutput.TableNames
}

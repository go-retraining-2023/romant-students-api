package data

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func CreateLocalClient(port int) *dynamodb.Client {

	endpointResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("http://localhost:%d", port),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile("local"),
		config.WithEndpointResolverWithOptions(endpointResolver),
	)

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

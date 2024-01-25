package data

import (
	"context"
	"log"

	awsAccess "github.com/RomanTykhyi/students-api/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func CreateLocalClient(dynamoDbUrl string, awsConfig *awsAccess.AwsConfig) (*dynamodb.Client, error) {
	log.Printf("Dynamo url:%v", dynamoDbUrl)
	endpointResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: dynamoDbUrl,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsConfig.AwsDefaultRegion),
		config.WithEndpointResolverWithOptions(endpointResolver),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     awsConfig.AwsAccessKey,
				SecretAccessKey: awsConfig.AwsSecretAccessKey,
				SessionToken:    awsConfig.AwsSessionToken,
				Source:          "Mock credentials used above for local instance",
			},
		}))

	if err != nil {
		return nil, err
	}

	return dynamodb.NewFromConfig(cfg), err
}

func CreateClient() (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile("roman-aws"),
	)

	if err != nil {
		return nil, err
	}

	return dynamodb.NewFromConfig(cfg), nil
}

func GetTables(dynamoDbClient *dynamodb.Client) ([]string, error) {
	tablesOutput, err := dynamoDbClient.ListTables(
		context.TODO(),
		&dynamodb.ListTablesInput{})

	if err != nil {
		return nil, err
	}

	return tablesOutput.TableNames, nil
}

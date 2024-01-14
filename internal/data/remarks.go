package data

import (
	"context"
	"fmt"

	"github.com/RomanTykhyi/students-api/internal/models"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

type StudentsRepository struct {
	dbClient  *dynamodb.Client
	tableName string
	partition string
}

func NewStudentsRepository(tableName, partition string, dbClient *dynamodb.Client) *StudentsRepository {
	return &StudentsRepository{
		dbClient:  dbClient,
		tableName: tableName,
		partition: partition,
	}
}

// реалізація методів інтерфейсу StudentsStore

func (repo *StudentsRepository) PutStudent(student *models.Student) error {
	student.PartitionId = repo.partition
	studentJson, err := attributevalue.MarshalMap(student)
	if err != nil {
		return err
	}

	fmt.Printf("%v", studentJson)

	putOutput, err := repo.dbClient.PutItem(
		context.TODO(),
		&dynamodb.PutItemInput{
			TableName:           &repo.tableName,
			Item:                studentJson,
			ConditionExpression: aws.String("attribute_not_exists(PK)"),
		},
	)

	if err != nil {
		return err
	}

	fmt.Printf("Successfully inserted: %v", putOutput)

	return nil
}


type StudentsRepository struct {
    dbClient *dynamodb.Client
    tableName string
    partition string
}

func NewStudentsRepository(tableName, partition string, dbClient *dynamodb.Client) *StudentsRepository {
    return &StudentsRepository{
        dbClient:   dbClient,
        tableName:  tableName,
        partition:  partition,
    }
}

// Решта методів будуть використовувати `repo.dbClient` замість `retrieveDynamoClient`
// Решта методів будуть використовувати repo.dbClient замість retrieveDynamoClient

//...

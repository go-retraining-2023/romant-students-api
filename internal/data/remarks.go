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

package data

import (
	"context"
	"fmt"
	"github.com/RomanTykhyi/students-api/internal/models"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"log"
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

func (repo *StudentsRepository) QueryStudents() []models.Student {
	var students []models.Student

	keyExpression := expression.Key("PartitionId").Equal(expression.Value(repo.partition))
	expr, err := expression.NewBuilder().WithKeyCondition(keyExpression).Build()
	if err != nil {
		log.Fatal("Error building students query.")
	}

	response, err := repo.dbClient.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:                 &repo.tableName,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})

	if err != nil {
		log.Printf("Error occured: %v", err)
		return nil
	}

	err = attributevalue.UnmarshalListOfMaps(response.Items, &students)
	if err != nil {
		log.Printf("Couldn't unmarshal query response. Here's why: %v\n", err)
	}

	return students
}

func (repo *StudentsRepository) GetStudent(id string) *models.Student {
	student := &models.Student{}
	key := map[string]types.AttributeValue{"PartitionId": &dynamodb.AttributeValue{S: &repo.partition}, "Id": &dynamodb.AttributeValue{S: &id}}

	output, err := repo.dbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key:       key,
		TableName: &repo.tableName,
	})

	if err != nil {
		return nil
	}

	err = attributevalue.UnmarshalMap(output.Item, student)
	if err != nil {
		log.Fatal("Error while deserializing student.")
	}

	return student
}

func (repo *StudentsRepository) UpdateStudent(student *models.Student) *models.Student {
	key := map[string]types.AttributeValue{"PartitionId": &dynamodb.AttributeValue{S: &repo.partition}, "Id": &dynamodb.AttributeValue{S: &student.Id}}
	update := expression.Set(expression.Name("FullName"), expression.Value(&student.FullName))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		log.Fatal("Error while creating the update student expression.")
	}

	output, err := repo.dbClient.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName:                 &repo.tableName,
		Key:                       key,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ReturnValues:              types.ReturnValueUpdatedNew,
	})

	if err != nil {
		log.Fatal("Error while updating student.")
	}

	err = attributevalue.UnmarshalMap(output.Attributes, student)
	if err != nil {
		log.Fatal("Couldn't unmarshall update student response.")
	}
	return student
}

func (repo *StudentsRepository) DeleteStudent(id string) bool {
	key := map[string]types.AttributeValue{"PartitionId": &dynamodb.AttributeValue{S: &repo.partition}, "Id": &dynamodb.AttributeValue{S: &id}}

	deleteOutput, err := repo.dbClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: &repo.tableName,
		Key:       key,
	})

	return err == nil && deleteOutput != nil
}

package data

import (
	"context"
	"fmt"
	"log"

	"github.com/RomanTykhyi/students-api/internal/di"
	"github.com/RomanTykhyi/students-api/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const TABLE_NAME = "Students"
const PARTITION = "students"

func retrieveDynamoClient() *dynamodb.Client {
	client, err := di.GetAppContainer().Resolve("dynamo-client")
	if err != nil {
		panic("Cannot get dynamo client")
	}
	return client.(*dynamodb.Client)
}

func getStudentKeyMap(studentId string) map[string]types.AttributeValue {
	partitionId, err := attributevalue.Marshal(PARTITION)
	if err != nil {
		log.Fatal(err)
	}

	id, err := attributevalue.Marshal(studentId)
	if err != nil {
		log.Fatal(err)
	}

	return map[string]types.AttributeValue{"PartitionId": partitionId, "Id": id}
}

type StudentsStore interface {
	PutStudent(student *models.Student) error
	QueryStudents() []models.Student
	GetStudent(id string) *models.Student
	UpdateStudent(student *models.Student) *models.Student
	DeleteStudent(id string) bool
}

type StudentsRepository struct{}

func (repo StudentsRepository) PutStudent(student *models.Student) error {
	dbClient := retrieveDynamoClient()

	student.PartitionId = PARTITION
	studentJson, err := attributevalue.MarshalMap(student)
	if err != nil {
		return err
	}

	fmt.Printf("%v", studentJson)

	putOutput, err := dbClient.PutItem(
		context.TODO(),
		&dynamodb.PutItemInput{
			TableName:           aws.String(TABLE_NAME),
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

func (repo StudentsRepository) QueryStudents() []models.Student {
	dbClient := retrieveDynamoClient()

	var response *dynamodb.QueryOutput
	var students []models.Student

	keyExpression := expression.Key("PartitionId").Equal(expression.Value(PARTITION))
	expr, err := expression.NewBuilder().WithKeyCondition(keyExpression).Build()
	if err != nil {
		log.Fatal("Error building students query.")
	}

	response, err = dbClient.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:                 aws.String(TABLE_NAME),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})

	if err != nil {
		log.Fatal("Error while querying students.")
	}

	err = attributevalue.UnmarshalListOfMaps(response.Items, &students)
	if err != nil {
		log.Printf("Couldn't unmarshal query response. Here's why: %v\n", err)
	}

	return students
}

func (repo StudentsRepository) GetStudent(id string) *models.Student {
	dbClient := retrieveDynamoClient()

	student := models.Student{}
	key := getStudentKeyMap(id)

	output, err := dbClient.GetItem(
		context.TODO(),
		&dynamodb.GetItemInput{
			Key:       key,
			TableName: aws.String(TABLE_NAME),
		})

	if err != nil {
		log.Fatal("Error while querying student.")
	}

	if output.Item == nil {
		return nil
	}

	err = attributevalue.UnmarshalMap(output.Item, &student)
	if err != nil {
		log.Fatal("Error while deserializing student.")
	}

	return &student
}

func (repo StudentsRepository) UpdateStudent(student *models.Student) *models.Student {
	dbClient := retrieveDynamoClient()

	key := getStudentKeyMap(student.Id)

	update := expression.Set(expression.Name("FullName"), expression.Value(student.FullName))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		log.Fatal("Error while creating the update student expression.")
	}

	output, err := dbClient.UpdateItem(
		context.TODO(),
		&dynamodb.UpdateItemInput{
			TableName:                 aws.String(TABLE_NAME),
			Key:                       key,
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
			ReturnValues:              types.ReturnValueUpdatedNew,
		})

	if err != nil {
		log.Fatal("Error while updating student.")
	}

	err = attributevalue.UnmarshalMap(output.Attributes, &student)
	if err != nil {
		log.Fatal("Couldn't unmarshall update student response.")
	}
	return student
}

func (repo StudentsRepository) DeleteStudent(id string) bool {
	dbClient := retrieveDynamoClient()

	key := getStudentKeyMap(id)

	deleteOutput, err := dbClient.DeleteItem(
		context.TODO(),
		&dynamodb.DeleteItemInput{
			TableName: aws.String(TABLE_NAME),
			Key:       key,
		})

	return err == nil && deleteOutput != nil
}

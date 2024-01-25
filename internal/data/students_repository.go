package data

import (
	"context"
	"fmt"
	"log"

	"github.com/RomanTykhyi/students-api/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func getStudentKeyMap(studentId string) (map[string]types.AttributeValue, error) {
	partitionId, err := attributevalue.Marshal(partitionId)
	if err != nil {
		return nil, err
	}

	id, err := attributevalue.Marshal(studentId)
	if err != nil {
		return nil, err
	}

	return map[string]types.AttributeValue{"PartitionId": partitionId, "Id": id}, nil
}

const (
	tableName   = "Students"
	partitionId = "students"
)

type StudentsStore interface {
	PutStudent(student *models.Student) error
	QueryStudents() ([]models.Student, error)
	GetStudent(id string) (*models.Student, error)
	UpdateStudent(student *models.Student) (*models.Student, error)
	DeleteStudent(id string) (bool, error)
}

type StudentsRepository struct {
	dbClient *dynamodb.Client
}

func NewStudentsRepository(dbClient *dynamodb.Client) *StudentsRepository {
	return &StudentsRepository{
		dbClient: dbClient,
	}
}

func (repo StudentsRepository) PutStudent(student *models.Student) error {
	student.PartitionId = partitionId

	studentJson, err := attributevalue.MarshalMap(student)
	if err != nil {
		return err
	}

	putOutput, err := repo.dbClient.PutItem(
		context.TODO(),
		&dynamodb.PutItemInput{
			TableName:           aws.String(tableName),
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

func (repo StudentsRepository) QueryStudents() ([]models.Student, error) {
	var response *dynamodb.QueryOutput
	var students []models.Student

	keyExpression := expression.Key("PartitionId").Equal(expression.Value(partitionId))
	expr, err := expression.NewBuilder().WithKeyCondition(keyExpression).Build()
	if err != nil {
		return nil, err
	}

	response, err = repo.dbClient.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})
	if err != nil {
		log.Printf("Error occured: %v", err)
		return nil, err
	}

	err = attributevalue.UnmarshalListOfMaps(response.Items, &students)
	if err != nil {
		return nil, err
	}

	return students, err
}

func (repo StudentsRepository) GetStudent(id string) (*models.Student, error) {
	student := models.Student{}
	key, err := getStudentKeyMap(id)
	if err != nil {
		return nil, err
	}

	output, err := repo.dbClient.GetItem(
		context.TODO(),
		&dynamodb.GetItemInput{
			Key:       key,
			TableName: aws.String(tableName),
		})
	if err != nil {
		return nil, err
	}

	err = attributevalue.UnmarshalMap(output.Item, &student)
	if err != nil {
		return nil, err
	}

	return &student, nil
}

func (repo StudentsRepository) UpdateStudent(student *models.Student) (*models.Student, error) {
	key, err := getStudentKeyMap(student.Id)
	if err != nil {
		return nil, err
	}

	update := expression.Set(expression.Name("FullName"), expression.Value(student.FullName))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return nil, err
	}

	output, err := repo.dbClient.UpdateItem(
		context.TODO(),
		&dynamodb.UpdateItemInput{
			TableName:                 aws.String(tableName),
			Key:                       key,
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
			ReturnValues:              types.ReturnValueUpdatedNew,
		})
	if err != nil {
		return nil, err
	}

	err = attributevalue.UnmarshalMap(output.Attributes, &student)
	if err != nil {
		return nil, err
	}
	return student, nil
}

func (repo StudentsRepository) DeleteStudent(id string) (bool, error) {
	key, err := getStudentKeyMap(id)
	if err != nil {
		return false, nil
	}

	_, err = repo.dbClient.DeleteItem(
		context.TODO(),
		&dynamodb.DeleteItemInput{
			TableName: aws.String(tableName),
			Key:       key,
		})
	if err != nil {
		return false, err
	}

	return true, nil
}

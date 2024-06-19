package repositories

import (
	"GameDay-API/internal/aws/awsclient"
	"GameDay-API/internal/models"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type ExportRepository interface {
	ListExports() ([]models.Export, error)
	GetExport(exportId string) (*models.Export, error)
	PutExport(export models.Export) error
}

type DynamoDbExportRepository struct {
	db        *dynamodb.Client
	tablename string
	ctx       context.Context
}

func NewExportRepository(context context.Context, tableName string) (ExportRepository, error) {
	db, err := awsclient.GetDynamodbClient()
	if err != nil {
		return nil, fmt.Errorf("error when initialising connection to DDB: %w", err)
	}

	return &DynamoDbExportRepository{
		db:        db,
		tablename: tableName,
		ctx:       context,
	}, nil
}

func (r *DynamoDbExportRepository) ListExports() ([]models.Export, error) {
	var response *dynamodb.ScanOutput
	var exports []models.Export
	response, err := r.db.Scan(r.ctx, &dynamodb.ScanInput{
		TableName: aws.String(r.tablename),
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't scan exports: %v", err)
	} else {
		err = attributevalue.UnmarshalListOfMaps(response.Items, &exports)
		if err != nil {
			return nil, fmt.Errorf("couldn't unmarshal query response: %v", err)
		}
	}
	return exports, nil
}

func (r *DynamoDbExportRepository) GetExport(exportId string) (*models.Export, error) {
	var response *dynamodb.QueryOutput
	var exports []models.Export
	keyEx := expression.Key("id").Equal(expression.Value(exportId))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return nil, fmt.Errorf("couldn't build expression for query: %w", err)
	}
	response, err = r.db.Query(r.ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.tablename),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't query exports for %v: %v", exportId, err)
	} else {
		err = attributevalue.UnmarshalListOfMaps(response.Items, &exports)
		if err != nil {
			return nil, fmt.Errorf("couldn't unmarshal query response: %v", err)
		}
	}
	if len(exports) == 0 {
		return nil, fmt.Errorf("couldn't find exactly one export matching %v: %v", exportId, err)
	}
	return &exports[0], nil
}

func (r *DynamoDbExportRepository) PutExport(export models.Export) error {
	item, err := attributevalue.MarshalMap(export)
	if err != nil {
		return fmt.Errorf("error when trying to convert export data to dynamodbattribute: %w", err)
	}
	params := &dynamodb.PutItemInput{
		TableName: aws.String(r.tablename),
		Item:      item,
	}
	if _, err := r.db.PutItem(r.ctx, params); err != nil {
		return fmt.Errorf("error when trying to persist: %w", err)
	}
	return nil
}

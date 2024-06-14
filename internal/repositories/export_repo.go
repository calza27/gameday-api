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
		return nil, fmt.Errorf("Error when initialising connection to DDB: %w", err)
	}

	return &DynamoDbExportRepository{
		db:        db,
		tablename: tableName,
		ctx:       context,
	}, nil
}

type ExportEntity struct {
	Id      string `dynamodbav:"id"`
	Name    string `dynamodbav:"name"`
	PdfFile string `dynamodbav:"pdfFile"`
	CsvFile string `dynamodbav:"csvFile"`
}

func (r *DynamoDbExportRepository) ListExports() ([]models.Export, error) {
	var response *dynamodb.ScanOutput
	var exports []ExportEntity
	response, err := r.db.Scan(r.ctx, &dynamodb.ScanInput{
		TableName: aws.String(r.tablename),
	})
	if err != nil {
		return nil, fmt.Errorf("Couldn't scan exports: %v\n", err)
	} else {
		err = attributevalue.UnmarshalListOfMaps(response.Items, &exports)
		if err != nil {
			return nil, fmt.Errorf("Couldn't unmarshal query response: %v\n", err)
		}
	}
	var exportModels []models.Export
	for _, export := range exports {
		exportModels = append(exportModels, convertToExportModel(export))
	}
	return exportModels, nil
}

func (r *DynamoDbExportRepository) GetExport(exportId string) (*models.Export, error) {
	var response *dynamodb.QueryOutput
	var exports []ExportEntity
	keyEx := expression.Key("id").Equal(expression.Value(exportId))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return nil, fmt.Errorf("Couldn't build expression for query: %w", err)
	}
	response, err = r.db.Query(r.ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.tablename),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})
	if err != nil {
		return nil, fmt.Errorf("Couldn't query exports for %v: %v\n", exportId, err)
	} else {
		err = attributevalue.UnmarshalListOfMaps(response.Items, &exports)
		if err != nil {
			return nil, fmt.Errorf("Couldn't unmarshal query response: %v\n", err)
		}
	}
	if len(exports) == 0 {
		return nil, fmt.Errorf("Couldn't find exactly one export matching %v: %v\n", exportId, err)
	}
	exportModel := convertToExportModel(exports[0])
	return &exportModel, nil
}

func (r *DynamoDbExportRepository) PutExport(export models.Export) error {
	exportEntity := convertToExportEntity(export)
	item, err := attributevalue.MarshalMap(exportEntity)
	if err != nil {
		return fmt.Errorf("Error when trying to convert export data to dynamodbattribute: %w", err)
	}
	params := &dynamodb.PutItemInput{
		TableName: aws.String(r.tablename),
		Item:      item,
	}
	if _, err := r.db.PutItem(r.ctx, params); err != nil {
		return fmt.Errorf("Error when trying to persist: %w", err)
	}
	return nil
}

func convertToExportModel(export ExportEntity) models.Export {
	return models.Export{
		Id:      export.Id,
		Name:    export.Name,
		PdfFile: export.PdfFile,
		CsvFile: export.CsvFile,
	}
}

func convertToExportEntity(export models.Export) ExportEntity {
	return ExportEntity{
		Id:      export.Id,
		Name:    export.Name,
		PdfFile: export.PdfFile,
		CsvFile: export.CsvFile,
	}
}

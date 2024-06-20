package main

import (
	"GameDay-API/internal/aws/awsclient"
	"GameDay-API/internal/filebuilder"
	"GameDay-API/internal/models"
	"GameDay-API/internal/repositories"
	"GameDay-API/internal/utils"
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/google/uuid"
)

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	reqBody := request.Body
	if reqBody == "" {
		return utils.BuildResponse("No request body supplied!", 400, nil), nil
	}
	var gameData models.GameData
	err := json.Unmarshal([]byte(reqBody), &gameData)
	if err != nil {
		return utils.BuildResponse("Error unmarshalling request body", 400, nil), nil
	}
	s3BucketName, err := GetS3BucketeName(ctx)
	if err != nil {
		return utils.BuildResponse(err.Error(), 500, nil), nil
	}
	exportsTable, err := GetDynamoDBTableName(ctx)
	if err != nil {
		return utils.BuildResponse(err.Error(), 500, nil), nil
	}
	s3Client, err := awsclient.GetS3Client()
	if err != nil {
		return utils.BuildResponse(err.Error(), 500, nil), nil
	}
	exportRepo, err := repositories.NewExportRepository(ctx, exportsTable)
	if err != nil {
		return utils.BuildResponse(err.Error(), 500, nil), nil
	}

	csvFileName := fmt.Sprintf("%s.csv", uuid.New().String())
	csvFile, err := filebuilder.BuildCsv(gameData)
	if err != nil {
		return utils.BuildResponse(err.Error(), 500, nil), nil
	}
	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(csvFileName),
		Body:   bytes.NewReader(csvFile),
	})
	if err != nil {
		return utils.BuildResponse(err.Error(), 500, nil), nil
	}

	pdfFileName := fmt.Sprintf("%s.pdf", uuid.New().String())
	pdfFile, err := filebuilder.BuildPdf(gameData)
	if err != nil {
		return utils.BuildResponse(err.Error(), 500, nil), nil
	}
	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(pdfFileName),
		Body:   bytes.NewReader(pdfFile),
	})
	if err != nil {
		return utils.BuildResponse(err.Error(), 500, nil), nil
	}

	export := models.Export{
		Id:      uuid.New().String(),
		Name:    buildRecordName(gameData),
		PdfFile: pdfFileName,
		CsvFile: csvFileName,
	}
	err = exportRepo.PutExport(export)
	if err != nil {
		return utils.BuildResponse(err.Error(), 500, nil), nil
	}
	return utils.BuildResponse("", 201, nil), nil
}

func buildRecordName(gameData models.GameData) string {
	return fmt.Sprintf("%s Vs %s @ %s, %s", gameData.TeamAAbbr, gameData.TeamBAbbr, gameData.Venue, gameData.GameDate)
}

func GetS3BucketeName(ctx context.Context) (string, error) {
	ssmClient, err := awsclient.GetSsmClient()
	if err != nil {
		return "", fmt.Errorf("error initializing connection to SSM: %w", err)
	}
	params := &ssm.GetParameterInput{
		Name: aws.String("/gameday-api-processor/s3/name"),
	}
	s3Name, err := ssmClient.GetParameter(ctx, params)
	if err != nil {
		return "", fmt.Errorf("error when getting parameter: %w", err)
	}
	return *s3Name.Parameter.Value, nil
}

func GetDynamoDBTableName(ctx context.Context) (string, error) {
	ssmClient, err := awsclient.GetSsmClient()
	if err != nil {
		return "", fmt.Errorf("error initializing connection to SSM: %w", err)
	}
	params := &ssm.GetParameterInput{
		Name: aws.String("/gameday-api-processor/ddb/name"),
	}
	ddbName, err := ssmClient.GetParameter(ctx, params)
	if err != nil {
		return "", fmt.Errorf("error when getting parameter: %w", err)
	}
	return *ddbName.Parameter.Value, nil
}

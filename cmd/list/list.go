package main

import (
	"GameDay-API/internal/aws/awsclient"
	"GameDay-API/internal/repositories"
	"GameDay-API/internal/utils"
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	exportsTable, err := GetDynamoDBTableName(ctx)
	if err != nil {
		utils.BuildResponse(err.Error(), 500, nil)
	}
	exportRepo, err := repositories.NewExportRepository(ctx, exportsTable)
	if err != nil {
		utils.BuildResponse(err.Error(), 500, nil)
	}
	exports, err := exportRepo.ListExports()
	if err != nil {
		utils.BuildResponse(err.Error(), 500, nil)
	}
	data, err := json.Marshal(exports)
	if err != nil {
		utils.BuildResponse(err.Error(), 500, nil)
	}
	return utils.BuildResponse(string(data), 200, nil), nil
}

func GetDynamoDBTableName(ctx context.Context) (string, error) {
	ssmClient, err := awsclient.GetSsmClient()
	if err != nil {
		return "", fmt.Errorf("Error initializing connection to SSM: %w", err)
	}
	params := &ssm.GetParameterInput{
		Name: aws.String("/gameday-api-processor/ddb/name"),
	}
	ddbName, err := ssmClient.GetParameter(ctx, params)
	if err != nil {
		return "", fmt.Errorf("Error when getting parameter: %w", err)
	}
	return *ddbName.Parameter.Value, nil
}

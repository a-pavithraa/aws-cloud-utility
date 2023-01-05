package service

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDbTableDetails struct {
	TableName   string `json:"tableName"`
	BillingType string `json:"billingType"`
	Region      string `json:"region"`
}

func (dd DynamoDbTableDetails) GetRegion() string {
	return dd.Region
}

type DynamoDBDescribeTableClient interface {
	ListTables(ctx context.Context, params *dynamodb.ListTablesInput, optFns ...func(options *dynamodb.Options)) (*dynamodb.ListTablesOutput, error)
	DescribeTable(ctx context.Context, params *dynamodb.DescribeTableInput, optFns ...func(options *dynamodb.Options)) (*dynamodb.DescribeTableOutput, error)
}

func NewDynamoDbClient(ctx context.Context, region string) (*dynamodb.Client, error) {

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %s", err)
	}

	dynamoDbClient := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.Region = region

	})

	return dynamoDbClient, nil

}

func GetDynamoDbTables(ctx context.Context, region string, client DynamoDBDescribeTableClient) ([]RestResponse, error) {
	var dynamboDBTablesList []RestResponse
	tables, err := client.ListTables(
		ctx, &dynamodb.ListTablesInput{})
	if err != nil {
		return nil, err
	} else {
		tableNames := tables.TableNames
		for _, tableName := range tableNames {
			dynamoDbTableDetails := DynamoDbTableDetails{
				TableName: tableName,
				Region:    region,
			}
			tableDetails, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
				TableName: &tableName,
			})
			if err != nil {
				return nil, err
			}
			billingMode := tableDetails.Table.BillingModeSummary
			if billingMode != nil {
				dynamoDbTableDetails.BillingType = string(billingMode.BillingMode)
			}

			dynamboDBTablesList = append(dynamboDBTablesList, dynamoDbTableDetails)
		}

	}
	return dynamboDBTablesList, nil
}

func ListDynamoDBTablesForRegion(ctx context.Context, region string) ([]RestResponse, error) {
	client, err := NewDynamoDbClient(ctx, region)
	if err != nil {
		return nil, err
	}
	return GetDynamoDbTables(ctx, region, client)

}

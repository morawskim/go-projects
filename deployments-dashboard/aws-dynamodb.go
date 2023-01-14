package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
)

type deployment struct {
	Project  string
	Instance string
	Version  string
}

func getDeployments(region string, tableName string) []*deployment {
	deployments := make([]*deployment, 0, 10)
	dynamodbSvc := dynamodb.New(session.New(), &aws.Config{Region: aws.String(region)})
	param := dynamodb.ScanInput{
		TableName: &tableName,
	}

	resp, err := dynamodbSvc.Scan(&param)
	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		log.Println(err.Error())
		return nil
	}

	for _, item := range resp.Items {
		if err := isValidDeploymentItem(item); err != nil {
			continue
		}

		deployments = append(deployments, &deployment{
			Project:  *item["Repository"].S,
			Instance: *item["Deployment"].S,
			Version:  *item["Version"].S,
		})
	}

	return deployments
}

func isValidDeploymentItem(item map[string]*dynamodb.AttributeValue) error {
	for _, attrName := range []string{"Repository", "Deployment", "Version"} {
		_, ok := item[attrName]
		if !ok {
			return fmt.Errorf("repository attribute %v not exists in item", attrName)
		}
	}

	return nil
}

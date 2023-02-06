package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"os"
)

type AppConfig struct {
	Repos []string
}

func main() {
	region := flag.String("region", "us-east-1", "AWS Region")
	accountID := flag.String("account", "", "Your AWS AccountID")
	config := flag.String("config", "", "The path to config file")
	tableName := flag.String("table", "", "The DynamoDB table name")
	flag.Parse()

	err := areFlagsValid(*region, *accountID, *config, *tableName)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		fmt.Println()
		os.Exit(1)
	}

	appConfig, err := parseConfigFile(*config)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		fmt.Println()
		os.Exit(1)
	}

	ecrHelper := ECRHelp{
		svc: &AwsEcrSvc{
			client: ecr.New(session.New(), &aws.Config{Region: aws.String(*region)}),
		},
	}

	for _, repo := range appConfig.Repos {
		fmt.Println(repo)
		images := ecrHelper.getAllAwsEcrImages(*accountID, repo)
		images = filterImagesToOnlySemVersion(images)
		sortImagesBySemVersion(images)
		for _, tag := range images {
			fmt.Println(*tag)
		}
	}

	fmt.Println("")
	fmt.Println("")
	fmt.Println("==== Deployments ====")
	fmt.Println("")

	deployments := getDeployments(*region, *tableName)
	for _, deploymentInfo := range deployments {
		fmt.Println(deploymentInfo.Instance, deploymentInfo.Project, deploymentInfo.Version)
	}
}

func areFlagsValid(region string, accountId string, config string, tableName string) error {
	if 0 == len(region) {
		return fmt.Errorf("you need to pass a region flag")
	}

	if 0 == len(accountId) {
		return fmt.Errorf("you need to pass a account flag")
	}

	if 0 == len(config) {
		return fmt.Errorf("you need to pass a config parametr")
	}

	_, err := os.Stat(config)

	if err != nil && os.IsNotExist(err) {
		return fmt.Errorf("the file passed in the config paramtere (%v) does not exist", config)
	}

	if err != nil {
		return fmt.Errorf("confg stat error: %w", err)
	}

	if 0 == len(tableName) {
		return fmt.Errorf("you need to pass a table flag")
	}

	return nil
}

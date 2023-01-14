package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	region := flag.String("region", "us-east-1", "AWS Region")
	accountID := flag.String("account", "", "Your AWS AccountID")
	repo := flag.String("repo", "", "The ECR Repo")
	tableName := flag.String("table", "", "The DynamoDB table name")
	flag.Parse()

	err := areFlagsValid(*region, *accountID, *repo, *tableName)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		fmt.Println()
		os.Exit(1)
	}

	images := getAllAwsEcrImages(*region, *accountID, *repo)
	images = filterImagesToOnlySemVersion(images)
	sortImagesBySemVersion(images)
	for _, tag := range images {
		fmt.Println(*tag)
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

func areFlagsValid(region string, accountId string, repo string, tableName string) error {
	if 0 == len(region) {
		return fmt.Errorf("you need to pass a region flag")
	}

	if 0 == len(accountId) {
		return fmt.Errorf("you need to pass a account flag")
	}

	if 0 == len(repo) {
		return fmt.Errorf("you need to pass a repo flag")
	}

	if 0 == len(tableName) {
		return fmt.Errorf("you need to pass a table flag")
	}

	return nil
}

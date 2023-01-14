package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"log"
)

func getAllAwsEcrImages(region string, registryId string, repositoryName string) []*string {
	done := false
	ecrSvc := ecr.New(session.New(), &aws.Config{Region: aws.String(region)})
	var imageIds []*string

	params := &ecr.ListImagesInput{
		RepositoryName: aws.String(repositoryName),
		MaxResults:     aws.Int64(100),
		RegistryId:     aws.String(registryId),
	}
	for !done {
		resp, err := ecrSvc.ListImages(params)

		if err != nil {
			log.Println(err.Error())
			return nil
		}

		for _, imageID := range resp.ImageIds {
			if nil == imageID.ImageTag {
				continue
			}

			imageIds = append(imageIds, imageID.ImageTag)
		}
		if resp.NextToken == nil {
			done = true
		} else {
			params.NextToken = resp.NextToken
		}
	}

	return imageIds
}

package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"log"
)

type ECRHelp struct {
	svc AwsEcrSvcInterface
}

func (a *ECRHelp) getAllAwsEcrImages(registryId string, repositoryName string) []*string {
	done := false
	var imageIds []*string

	params := &ecr.ListImagesInput{
		RepositoryName: aws.String(repositoryName),
		MaxResults:     aws.Int64(100),
		RegistryId:     aws.String(registryId),
	}
	for !done {
		resp, err := a.svc.GetAllRepositoryImagesTags(params)

		if err != nil {
			log.Printf("Cannot get AWS ECR image rags: %v\n", err.Error())
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

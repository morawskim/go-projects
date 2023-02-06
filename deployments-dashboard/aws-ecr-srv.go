package main

import (
	"github.com/aws/aws-sdk-go/service/ecr"
)

type AwsEcrSvc struct {
	client *ecr.ECR
}

type MockedAwsEcrSvc struct {
	data *ecr.ListImagesOutput
}

type AwsEcrSvcInterface interface {
	GetAllRepositoryImagesTags(params *ecr.ListImagesInput) (*ecr.ListImagesOutput, error)
}

func (a *AwsEcrSvc) GetAllRepositoryImagesTags(params *ecr.ListImagesInput) (*ecr.ListImagesOutput, error) {
	return a.client.ListImages(params)
}

func (m *MockedAwsEcrSvc) GetAllRepositoryImagesTags(params *ecr.ListImagesInput) (*ecr.ListImagesOutput, error) {
	return m.data, nil
}

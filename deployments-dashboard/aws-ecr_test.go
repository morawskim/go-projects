package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetImages(t *testing.T) {
	ecrHelper := ECRHelp{
		svc: &MockedAwsEcrSvc{
			data: &ecr.ListImagesOutput{
				ImageIds: []*ecr.ImageIdentifier{
					{
						ImageDigest: nil,
						ImageTag:    nil,
					},
					{
						ImageDigest: aws.String("1234567890"),
						ImageTag:    aws.String("0.1.0"),
					},
					{
						ImageDigest: aws.String("abcdef01234567890"),
						ImageTag:    aws.String("0.2.0"),
					},
					{
						ImageDigest: aws.String("abcdef01234567890"),
						ImageTag:    aws.String("feature-foo"),
					},
				},
				NextToken: nil,
			},
		},
	}

	slice := ecrHelper.getAllAwsEcrImages("12345678901", "foo/bar")
	assert.Len(t, slice, 3)
}

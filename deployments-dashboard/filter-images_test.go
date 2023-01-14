package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testImagesData = []struct {
	image           string
	validSemVersion bool
}{
	{image: "FOO-123", validSemVersion: false},
	{image: "develop-520", validSemVersion: false},
	{image: "main", validSemVersion: false},
	{image: "stable", validSemVersion: false},
	{image: "2.30.29-beta", validSemVersion: false},
	{image: "2.20.1", validSemVersion: true},
}

func TestFilterImagesOnlySemVersion(t *testing.T) {
	images := []*string{
		&testImagesData[0].image,
		&testImagesData[1].image,
		&testImagesData[2].image,
		&testImagesData[3].image,
		&testImagesData[4].image,
		&testImagesData[5].image,
	}
	filteredImages := filterImagesToOnlySemVersion(images)

	assert.Len(t, filteredImages, 1)
	assert.Equal(t, "2.20.1", *filteredImages[0])
}

func TestIsValidSemVersion(t *testing.T) {
	for _, imageTestItem := range testImagesData {
		result := isValidSemVersion(imageTestItem.image)
		if imageTestItem.validSemVersion {
			assert.Nil(t, result)
		} else {
			assert.Error(t, result)
		}
	}
}

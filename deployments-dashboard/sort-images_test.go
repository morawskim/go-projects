package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSortImagesSemVersionDesc(t *testing.T) {
	tmp := []string{
		"1.12.3",
		"1.10.4",
		"1.10.0",
		"1.15.1",
		"1.12.4",
	}

	images := []*string{
		&tmp[0],
		&tmp[1],
		&tmp[2],
		&tmp[3],
		&tmp[4],
	}

	sortImagesBySemVersion(images)
	assert.Equal(t, "1.15.1", *images[0])
	assert.Equal(t, "1.12.4", *images[1])
	assert.Equal(t, "1.12.3", *images[2])
	assert.Equal(t, "1.10.4", *images[3])
	assert.Equal(t, "1.10.0", *images[4])
}

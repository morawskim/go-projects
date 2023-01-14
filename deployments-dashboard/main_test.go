package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAreFlagsValid(t *testing.T) {
	var err error
	err = areFlagsValid("", "", "", "")
	assert.ErrorContains(t, err, "region")

	err = areFlagsValid("foo", "", "", "")
	assert.ErrorContains(t, err, "account")

	err = areFlagsValid("foo", "12345678901", "", "")
	assert.ErrorContains(t, err, "repo")

	err = areFlagsValid("foo", "12345678901", "foo", "")
	assert.ErrorContains(t, err, "table")

	err = areFlagsValid("foo", "12345678901", "foo", "deployments")
	assert.Nil(t, err)
}

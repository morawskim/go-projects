package main

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
)

func filterImagesToOnlySemVersion(images []*string) []*string {
	filteredSlice := make([]*string, 0, len(images)/2)

	for _, image := range images {
		err := isValidSemVersion(*image)
		if err != nil {
			continue
		}

		filteredSlice = append(filteredSlice, image)
	}

	return filteredSlice
}

func isValidSemVersion(semVersion string) error {
	version, err := semver.NewVersion(semVersion)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Checking %v", semVersion))
	}

	if version.Prerelease() != "" {
		return fmt.Errorf("version %v looks like prelease version", semVersion)
	}

	return nil
}

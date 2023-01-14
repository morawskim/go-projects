package main

import (
	"github.com/nsd20463/numeralsort"
	"sort"
)

func sortImagesBySemVersion(images []*string) {
	sort.Slice(images, func(i, j int) bool {
		return !numeralsort.Less(*images[i], *images[j])
	})
}

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func main() {
	kipReader := &kio.ByteReader{Reader: os.Stdin}
	err := kio.Pipeline{
		Inputs: []kio.Reader{kipReader},
		Filters: []kio.Filter{
			kio.FilterFunc(createFilterFunction(kipReader)),
		},
		Outputs:               []kio.Writer{kio.ByteWriter{Writer: os.Stdout}},
		ContinueOnEmptyResult: false,
	}.Execute()

	if err != nil {
		log.Fatal(err)
	}
}

func createFilterFunction(kipReader *kio.ByteReader) func(nodes []*yaml.RNode) ([]*yaml.RNode, error) {
	return func(nodes []*yaml.RNode) ([]*yaml.RNode, error) {
		cronJobsToDisable, cronJobsToEnable, err := extractCronJobsToDisableAndEnable(kipReader.FunctionConfig)
		if err != nil {
			return nil, err
		}

		filters := make([]yaml.Filter, 3)
		filters = []yaml.Filter{
			yaml.Tee(yaml.SetAnnotation("foo", "bar")),
			yaml.Lookup("spec"),
			nil,
		}

		for i := range nodes {
			resource := nodes[i]

			if "CronJob" == resource.GetKind() && "batch/v1" == resource.GetApiVersion() {
				var err error = nil
				if contains(cronJobsToDisable, resource.GetName()) {
					filters[2] = yaml.Tee(yaml.SetField("suspend", yaml.MustParse("true")))
					err = resource.PipeE(filters...)
				} else if contains(cronJobsToEnable, resource.GetName()) {
					filters[2] = yaml.Tee(yaml.SetField("suspend", yaml.MustParse("false")))
					err = resource.PipeE(filters...)
				}

				if err != nil {
					return nil, err
				}
			}
		}
		return nodes, nil
	}
}

func extractCronJobsToDisableAndEnable(y *yaml.RNode) ([]string, []string, error) {

	f := func(y *yaml.RNode, path string) ([]string, error) {
		value, err := y.GetFieldValue(path)
		if err != nil {
			var targetError yaml.NoFieldError
			if errors.As(err, &targetError) {
				return []string{}, nil
			} else {
				return nil, err
			}
		}

		if _, ok := value.([]interface{}); !ok {
			return nil, fmt.Errorf("expected []string for path %q, got %T", path, value)
		}

		slice := make([]string, 0, len(value.([]interface{})))

		for _, val := range value.([]interface{}) {
			if _, ok := val.(string); !ok {
				return nil, fmt.Errorf("expected slice of strings for path %q, got %T", path, val)
			}

			slice = append(slice, val.(string))
		}

		return slice, nil
	}

	cronJobsToDisable, err := f(y, "spec.cronJobsToDisable")
	if err != nil {
		return nil, nil, err
	}
	cronJobsToEnable, err := f(y, "spec.cronJobsToEnable")
	if err != nil {
		return nil, nil, err
	}

	return cronJobsToDisable, cronJobsToEnable, nil
}

func contains[T comparable](slice []T, search T) bool {
	for _, value := range slice {
		if value == search {
			return true
		}
	}
	return false
}

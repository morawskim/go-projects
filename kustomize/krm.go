package main

import (
	"log"
	"os"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func main() {
	err := kio.Pipeline{
		Inputs: []kio.Reader{&kio.ByteReader{Reader: os.Stdin}},
		Filters: []kio.Filter{
			kio.FilterFunc(func(nodes []*yaml.RNode) ([]*yaml.RNode, error) {
				for i := range nodes {
					resource := nodes[i]

					if "CronJob" == resource.GetKind() && "batch/v1" == resource.GetApiVersion() {
						err := resource.PipeE(
							yaml.Tee(yaml.SetAnnotation("foo", "bar")),
							yaml.Lookup("spec"),
							yaml.Tee(yaml.SetField("suspend", yaml.MustParse("true"))),
						)
						if err != nil {
							return nil, err
						}
					}
				}
				return nodes, nil
			}),
		},
		Outputs:               []kio.Writer{kio.ByteWriter{Writer: os.Stdout}},
		ContinueOnEmptyResult: false,
	}.Execute()

	if err != nil {
		log.Fatal(err)
	}
}

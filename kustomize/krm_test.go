package main

import (
	"os"
	"reflect"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"strings"
	"testing"
)

func mustGetYamlFileContent(path string) string {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(file)
}

func Test_extractCronJobsToDisableAndEnable(t *testing.T) {
	type args struct {
		yamlFilePath string
	}
	tests := []struct {
		name       string
		args       args
		want       []string
		want1      []string
		wantErr    bool
		errMessage string
	}{
		{
			name: "missing configuration fields",
			args: args{
				yamlFilePath: "./tests/missing-plugin-configuration.yml",
			},
			want:    []string{},
			want1:   []string{},
			wantErr: false,
		},
		{
			name: "only enable cronjob configuration",
			args: args{
				yamlFilePath: "./tests/plugin-configuration-only-for-enable.yml",
			},
			want:    []string{},
			want1:   []string{"my-cronjob-to-enable"},
			wantErr: false,
		},
		{
			name: "only disable cronjob configuration",
			args: args{
				yamlFilePath: "./tests/plugin-configuration-only-for-disable.yml",
			},
			want:    []string{"my-cronjob-to-suspend"},
			want1:   []string{},
			wantErr: false,
		},
		{
			name: "full configuration for plugin",
			args: args{
				yamlFilePath: "./tests/plugin-configuration-full.yml",
			},
			want:    []string{"my-cronjob-to-suspend"},
			want1:   []string{"my-cronjob-to-enable", "my-cronjob-to-enable2"},
			wantErr: false,
		},
		{
			name: "invalid type - map instead of list",
			args: args{
				yamlFilePath: "./tests/plugin-configuration-invalid-type.yml",
			},
			want:       nil,
			want1:      nil,
			wantErr:    true,
			errMessage: `expected []string for path "spec.cronJobsToDisable", got map[string]interface {}`,
		},
		{
			name: "invalid type - slice of int instead of string",
			args: args{
				yamlFilePath: "./tests/plugin-configuration-invalid-type-in-list.yml",
			},
			want:       nil,
			want1:      nil,
			wantErr:    true,
			errMessage: `expected slice of strings for path "spec.cronJobsToEnable", got int`,
		},
		{
			name: "invalid type - mixed types in slice",
			args: args{
				yamlFilePath: "./tests/plugin-configuration-mixed-types.yml",
			},
			want:       nil,
			want1:      nil,
			wantErr:    true,
			errMessage: `expected slice of strings for path "spec.cronJobsToEnable", got int`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := kio.ByteReader{
				Reader: strings.NewReader(mustGetYamlFileContent(tt.args.yamlFilePath)),
			}
			_, err := reader.Read()

			if err != nil {
				t.Errorf("cannot read a yaml file error = %v", err)
			}

			got, got1, err := extractCronJobsToDisableAndEnable(reader.FunctionConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractCronJobsToDisableAndEnable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errMessage {
				t.Errorf("extractCronJobsToDisableAndEnable() error = %v, wantErr %v", err.Error(), tt.errMessage)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractCronJobsToDisableAndEnable() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("extractCronJobsToDisableAndEnable() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

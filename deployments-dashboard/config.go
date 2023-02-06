package main

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"os"
)

func parseConfigFile(config string) (*AppConfig, error) {
	appConfig := new(AppConfig)
	configBytes, err := os.ReadFile(config)

	if err != nil {
		return nil, errors.Errorf("cannot read config file %v: %v", config, err)
	}

	err = yaml.Unmarshal(configBytes, appConfig)
	if err != nil {
		return nil, errors.Errorf("cannot parse config file %v: %v", config, err)
	}

	return appConfig, nil
}

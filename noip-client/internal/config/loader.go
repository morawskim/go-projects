package config

import "os"

type NoIpConfig struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
	Hostname string `validate:"required,fqdn"`
}

func CreateNoIpConfigFromEnvVariables() *NoIpConfig {
	return &NoIpConfig{
		Username: os.Getenv("NOIP_USERNAME"),
		Password: os.Getenv("NOIP_PASSWORD"),
		Hostname: os.Getenv("NOIP_HOSTNAME"),
	}
}

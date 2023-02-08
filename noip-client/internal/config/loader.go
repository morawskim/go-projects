package config

import "os"

type NoIpConfig struct {
	Username string
	Password string
	Hostname string
}

func CreateNoIpConfigFromEnvVariables() *NoIpConfig {
	return &NoIpConfig{
		Username: os.Getenv("NOIP_USERNAME"),
		Password: os.Getenv("NOIP_PASSWORD"),
		Hostname: os.Getenv("NOIP_HOSTNAME"),
	}
}

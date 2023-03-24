package config

import (
	"github.com/spf13/viper"
)

type NoIpConfig struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
	Hostname string `validate:"required,fqdn"`
}

func CreateFromViper() *NoIpConfig {
	return &NoIpConfig{
		Username: viper.GetString("NOIP_USERNAME"),
		Password: viper.GetString("NOIP_PASSWORD"),
		Hostname: viper.GetString("NOIP_HOSTNAME"),
	}
}

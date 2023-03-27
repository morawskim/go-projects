package cmd

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
	"noip-client/internal"
	"noip-client/internal/config"
	"noip-client/internal/iphelper"
	"os"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the current IP address assigned to noip domain",
	Run: func(cmd *cobra.Command, args []string) {
		logger := slog.New(slog.NewJSONHandler(os.Stdout))
		logger.Enabled(slog.LevelInfo)

		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			logger.Error("Unable to parse config file", err)
			os.Exit(1)
		}

		noIpConfig := config.CreateFromViper()

		v := validator.New()
		err := v.Struct(noIpConfig)
		if err != nil {
			logger.Error("Missing or invalid evn variables:", err)
			os.Exit(1)
		}

		ipHelper := iphelper.NewIpHelper()
		updateError := internal.UpdateNoIpDnsRecord(noIpConfig, ipHelper)
		if nil != updateError {
			if errors.Is(updateError, internal.ErrTheSameIpAddr) {
				logger.Info(updateError.Error())
			} else {
				logger.Error("Unable to update DNS record", updateError, slog.String("hostname", noIpConfig.Hostname), slog.String("username", noIpConfig.Username))
				os.Exit(1)
			}
		}
	},
}

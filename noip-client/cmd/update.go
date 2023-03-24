package cmd

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
	"noip-client/internal/config"
	"noip-client/internal/iphelper"
	"noip-client/internal/noip"
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
		currentAssignedIp, err := ipHelper.GetCurrentAssignedIp(noIpConfig.Hostname)
		if err != nil {
			logger.Error("Unable to get current assigned IP address", err, slog.String("hostname", noIpConfig.Hostname))
			os.Exit(1)
		}

		myCurrentPublicIp, err := ipHelper.GetCurrentPublicIpAddress()
		if err != nil {
			logger.Error("Unable to get current public IP address", err)
			os.Exit(1)
		}

		if !currentAssignedIp.Equal(myCurrentPublicIp) {
			noipApiClient := noip.NewApiClient(noIpConfig)
			updateApiErr := noipApiClient.UpdateAssignedIp(currentAssignedIp)
			if updateApiErr != nil {
				logger.Error("Unable to update assigned noip address", updateApiErr, slog.String("hostname", noIpConfig.Hostname), slog.String("username", noIpConfig.Username))
				os.Exit(1)
			}
		} else {
			logger.Info("Current assigned IP address is the same as public address", slog.String("noip", currentAssignedIp.String()), slog.String("publicIp", myCurrentPublicIp.String()))
		}
	},
}

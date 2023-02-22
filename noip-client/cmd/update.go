package cmd

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
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
		noIpConfig := config.CreateNoIpConfigFromEnvVariables()

		v := validator.New()
		err := v.Struct(noIpConfig)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		ipHelper := iphelper.NewIpHelper()
		currentAssignedIp, err := ipHelper.GetCurrentAssignedIp(noIpConfig.Hostname)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		myCurrentPublicIp, err := ipHelper.GetCurrentPublicIpAddress()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if !currentAssignedIp.Equal(myCurrentPublicIp) {
			noipApiClient := noip.NewApiClient(noIpConfig)
			updateApiErr := noipApiClient.UpdateAssignedIp(currentAssignedIp)
			if updateApiErr != nil {
				fmt.Println(updateApiErr)
				os.Exit(1)
			}
		}
	},
}

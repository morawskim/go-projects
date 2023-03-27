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
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func init() {
	rootCmd.AddCommand(daemonCmd)
}

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run forevwer and update ",
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

		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			c := time.Tick(5 * time.Minute)
			termChan := make(chan os.Signal)
			intChan := make(chan os.Signal)
			signal.Notify(termChan, syscall.SIGTERM)
			signal.Notify(intChan, syscall.SIGINT)

			defer close(termChan)
			defer close(intChan)

			for {
				select {
				case <-termChan:
					logger.Info("Received SIGTERM, exiting")
					wg.Done()
					return
				case <-intChan:
					logger.Info("Received SIGINT, exiting")
					wg.Done()
					return
				case <-c:
					updateError := internal.UpdateNoIpDnsRecord(noIpConfig, ipHelper)
					if nil != updateError {
						if errors.Is(updateError, internal.ErrTheSameIpAddr) {
							logger.Info(updateError.Error())
						} else {
							logger.Error("Unable to update DNS record", updateError, slog.String("hostname", noIpConfig.Hostname), slog.String("username", noIpConfig.Username))
						}
					}
				}
			}
		}()

		wg.Wait()
	},
}

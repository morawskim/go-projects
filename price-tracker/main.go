package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type item2 struct {
	productName string
	productUrl  string
}

type trackerConfig struct {
	mu        sync.RWMutex
	products  []item2
	selectors map[string]selector
	pr        map[string]string
}

func (t *trackerConfig) setData(products []item2, selectors map[string]selector, pr map[string]string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.products = products
	t.selectors = selectors
	t.pr = pr
}

func (t *trackerConfig) getData() ([]item2, map[string]selector, map[string]string) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.products, t.selectors, t.pr
}

type expressionEnv struct {
	Element *colly.HTMLElement
}

const productCtxKey string = "product"

func (expressionEnv) GetTextContent(el *colly.HTMLElement) string {
	return el.Text
}

func (expressionEnv) GetFirstChildTextContent(el *colly.HTMLElement, selectors ...string) string {
	for _, s := range selectors {
		find := el.DOM.Find(s)

		if find.Length() == 0 {
			continue
		}

		return find.First().Text()
	}

	return ""
}

func (expressionEnv) GetAttribute(el *colly.HTMLElement, attributeName string) string {
	return el.Attr(attributeName)
}

func (expressionEnv) GetInputValue(el *colly.HTMLElement) string {
	return el.Attr("value")
}

func reloadTrackerConfig(cfgFile string, t *trackerConfig) error {
	i := config{}
	b, err := loadConfigFile(cfgFile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(b, &i)
	if err != nil {
		return err
	}
	pc, mapPr := processConfig(&i)
	t.setData(pc, i.Selectors, mapPr)
	slog.Default().Info("Configuration reloaded", slog.String("file", cfgFile))
	return nil
}

var onlyDigitsRegex = regexp.MustCompile(`[^0-9.,]+`)

func main() {
	var cfgFile string
	var interval time.Duration

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	slog.SetDefault(logger)

	// Initialize a new Cobra command
	var rootCmd = &cobra.Command{
		Use:   "hello",
		Short: "Prints 'Hello, World!'",
		Run: func(cmd *cobra.Command, args []string) {
			if cfgFile == "" {
				cobra.CheckErr(fmt.Errorf("no config file specified"))
			}

			trackerConfigStruct := &trackerConfig{}
			err := reloadTrackerConfig(cfgFile, trackerConfigStruct)
			cobra.CheckErr(err)

			collectorMinPrice := newMinPriceCollector()
			ch := createChannel(collectorMinPrice)

			products, _, _ := trackerConfigStruct.getData()
			registerMetrics(products, collectorMinPrice)

			// Handle signals for SIGHUP
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, syscall.SIGHUP)

			// Ticker for 1h reload
			reloadTicker := time.NewTicker(1 * time.Hour)

			go func() {
				for {
					select {
					case sig := <-sigCh:
						slog.Default().Info("Received signal", slog.String("signal", sig.String()))
						if err := reloadTrackerConfig(cfgFile, trackerConfigStruct); err != nil {
							slog.Default().Error("Failed to reload config on signal", slog.String("error", err.Error()))
						}
					case <-reloadTicker.C:
						slog.Default().Info("Periodic config reload")
						if err := reloadTrackerConfig(cfgFile, trackerConfigStruct); err != nil {
							slog.Default().Error("Failed to periodic reload config", slog.String("error", err.Error()))
						}
					}
				}
			}()

			go runPeriodically(interval, trackerConfigStruct, ch)

			slog.Default().Info("starting http server")
			register(trackerConfigStruct, collectorMinPrice)
			close(ch)
		},
	}
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	rootCmd.PersistentFlags().DurationVar(&interval, "interval", 18*time.Hour, "interval")

	// Execute the command
	if err := rootCmd.Execute(); err != nil {
		slog.Default().Error(err.Error())
		os.Exit(1)
	}
}

func runPeriodically(interval time.Duration, t *trackerConfig, ch chan metric) {
	// Run the function immediately
	products, selectors, pr := t.getData()
	collect(products, selectors, pr, ch)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			products, selectors, pr := t.getData()
			collect(products, selectors, pr, ch)
		}
	}
}

func loadConfigFile(cfgFile string) ([]byte, error) {
	u, err := url.Parse(cfgFile)

	if err == nil && (u.Scheme == "http" || u.Scheme == "https") {
		response, err := http.Get(cfgFile)
		if err != nil {
			return nil, err
		}

		defer response.Body.Close()

		return io.ReadAll(response.Body)
	}

	if !isValidFile(cfgFile) {
		cobra.CheckErr(fmt.Errorf(`config file "%s" not exists`, cfgFile))
	}

	return os.ReadFile(cfgFile)
}

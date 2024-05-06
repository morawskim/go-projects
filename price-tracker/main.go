package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"regexp"
	"time"
)

type item2 struct {
	productName string
	productUrl  string
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

			if !isValidFile(cfgFile) {
				cobra.CheckErr(fmt.Errorf(`config file "%s" not exists`, cfgFile))
			}

			i := config{}
			b, err := os.ReadFile(cfgFile)
			cobra.CheckErr(err)

			err = yaml.Unmarshal(b, &i)
			cobra.CheckErr(err)
			pc, mapPr := processConfig(&i)

			ch := createChannel()
			registerMetrics(pc)
			go runPeriodically(interval, pc, i.Selectors, mapPr, ch)

			slog.Default().Info("starting http server")
			register()
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

func runPeriodically(interval time.Duration, products []item2, selectors map[string]selector, pr map[string]string, ch chan metric) {
	// Run the function immediately
	collect(products, selectors, pr, ch)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			collect(products, selectors, pr, ch)
		}
	}
}

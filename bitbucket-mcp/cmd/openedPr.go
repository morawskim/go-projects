package cmd

import (
	"bitbucket-mcp/bitbucket"
	"fmt"

	"github.com/spf13/cobra"
)

// openedPrCmd represents the openedPr command
var openedPrCmd = &cobra.Command{
	Use:   "openedPr",
	Short: "Get a list of opened pull requests",
	Run: func(cmd *cobra.Command, args []string) {
		pr := bitbucket.GetOpenPr()
		for _, p := range pr {
			fmt.Println(p)
		}
	},
}

func init() {
	rootCmd.AddCommand(openedPrCmd)
}

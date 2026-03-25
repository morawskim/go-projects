package cmd

import (
	"bitbucket-mcp/bitbucket"
	"context"
	"log"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

type Input struct {
}

type Output struct {
	OpenPullRequests string `json:"openPullRequests" jsonschema:"the opened pull requests which are waiting for review. each line contains the title of PR and link to PR"`
}

func GetOpenedPullRequests(ctx context.Context, req *mcp.CallToolRequest, input Input) (
	*mcp.CallToolResult,
	Output,
	error,
) {
	return nil, Output{OpenPullRequests: strings.Join(bitbucket.GetOpenPr(), "\n")}, nil
}

// mcpCmd represents the mcp command
var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Run Bitbucket MCP Server",
	Run: func(cmd *cobra.Command, args []string) {
		server := mcp.NewServer(&mcp.Implementation{Name: "bitbucket", Version: "v1.0.0"}, nil)
		mcp.AddTool(server, &mcp.Tool{Name: "openedPullRequests", Description: "Return opened pull requests in bitbucket which are waiting for code review. Each line contains the tile of pull requests and link to pull request"}, GetOpenedPullRequests)
		// Run the server over stdin/stdout, until the client disconnects.
		if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}

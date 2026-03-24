package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type config struct {
	bitbucketUsername     string
	bitbucketAppToken     string
	bitbucketRepo         string
	bitbucketPrAuthorUuid string
}

func (c *config) getRepoOwner() string {
	parts := strings.SplitN(c.bitbucketRepo, "/", 2)
	return parts[0]
}
func (c *config) getRepoSlug() string {
	i := strings.Index(c.bitbucketRepo, "/")
	if i == -1 {
		return ""
	}
	return c.bitbucketRepo[i+1:]
}

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
	return nil, Output{OpenPullRequests: strings.Join(getOpenPr(), "\n")}, nil
}

func main() {
	server := mcp.NewServer(&mcp.Implementation{Name: "bitbucket", Version: "v1.0.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "openedPullRequests", Description: "Return opened pull requests in bitbucket which are waiting for code review. Each line contains the tile of pull requests and link to pull request"}, GetOpenedPullRequests)
	// Run the server over stdin/stdout, until the client disconnects.
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}

func getConfiguration() *config {
	pflag.String("bitbucket-username", "", "Bitbucket username parameter")
	pflag.String("bitbucket-token", "", "Bitbucket App token")
	pflag.String("bitbucket-repo", "", "Bitbucket repo owner/repoSlug")
	pflag.String("bitbucket-pr-author-uuid", "", "Bitbucket PR author uuid")
	pflag.Parse()

	viper.BindPFlag("BITBUCKET_USERNAME", pflag.Lookup("bitbucket-username"))
	viper.BindPFlag("BITBUCKET_APP_TOKEN", pflag.Lookup("bitbucket-token"))
	viper.BindPFlag("BITBUCKET_PR_REPO", pflag.Lookup("bitbucket-repo"))
	viper.BindPFlag("BITBUCKET_PR_AUTHOR_UUID", pflag.Lookup("bitbucket-pr-author-uuid"))
	viper.AutomaticEnv()

	return &config{
		bitbucketUsername:     viper.GetString("BITBUCKET_USERNAME"),
		bitbucketAppToken:     viper.GetString("BITBUCKET_APP_TOKEN"),
		bitbucketRepo:         viper.GetString("BITBUCKET_PR_REPO"),
		bitbucketPrAuthorUuid: viper.GetString("BITBUCKET_PR_AUTHOR_UUID"),
	}
}

func getOpenPr() []string {
	cfg := getConfiguration()
	client, err := bitbucket.NewBasicAuth(cfg.bitbucketUsername, cfg.bitbucketAppToken)
	if err != nil {
		panic(err)
	}

	opt := &bitbucket.PullRequestsOptions{
		Owner:    cfg.getRepoOwner(),
		RepoSlug: cfg.getRepoSlug(),
		Query:    fmt.Sprintf(`author.uuid="{%s}" AND state="OPEN" AND draft=false`, cfg.bitbucketPrAuthorUuid),
	}
	res, err := client.Repositories.PullRequests.List(opt)
	if err != nil {
		panic(err)
	}

	resultMap := res.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	result := make([]string, len(values))

	for _, value := range values {
		pr := value.(map[string]interface{})
		result = append(result, fmt.Sprintf(
			"%s - https://bitbucket.org/%s/pull-requests/%d",
			pr["title"],
			cfg.bitbucketRepo,
			int(pr["id"].(float64))))
	}

	return result
}

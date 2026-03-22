package main

import (
	"fmt"
	"strings"

	"github.com/ktrysmt/go-bitbucket"
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

func main() {
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

	for _, value := range values {
		pr := value.(map[string]interface{})
		fmt.Println(fmt.Sprintf(
			"%s - https://bitbucket.org/%s/pull-requests/%d",
			pr["title"],
			cfg.bitbucketRepo,
			int(pr["id"].(float64))))
		fmt.Println()
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

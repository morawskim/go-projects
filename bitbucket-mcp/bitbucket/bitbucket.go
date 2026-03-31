package bitbucket

import (
	"fmt"
	"strconv"

	"github.com/ktrysmt/go-bitbucket"
)

func GetOpenPr() []string {
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
	result := make([]string, 0, len(values))

	for _, value := range values {
		pr := value.(map[string]interface{})
		prID := int(pr["id"].(float64))
		idStr := strconv.Itoa(prID)
		resDetails, err := client.Repositories.PullRequests.Get(&bitbucket.PullRequestsOptions{
			ID:       idStr,
			Owner:    cfg.getRepoOwner(),
			RepoSlug: cfg.getRepoSlug(),
		})
		if err != nil {
			panic(err)
		}
		resDetailsMap := resDetails.(map[string]interface{})
		participants := resDetailsMap["participants"].([]interface{})
		isApproved := false
		for _, participant := range participants {
			participantMap := participant.(map[string]interface{})
			if participantMap["approved"].(bool) {
				isApproved = true
				break
			}
		}

		if isApproved {
			continue
		}

		result = append(result, fmt.Sprintf(
			"%s - https://bitbucket.org/%s/pull-requests/%d",
			pr["title"],
			cfg.bitbucketRepo,
			prID))
	}

	return result
}

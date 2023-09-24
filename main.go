package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/v55/github"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Println("GITHUB_TOKEN environment variable not set")
		return
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opt := &github.ListOptions{}
	for {
		orgs, resp, err := client.Organizations.List(ctx, "", opt)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		for _, org := range orgs {
			repoOpt := &github.RepositoryListByOrgOptions{ListOptions: github.ListOptions{}}
			for {
				repos, resp, err := client.Repositories.ListByOrg(ctx, org.GetLogin(), repoOpt)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}

				for _, repo := range repos {
					if repo.GetPushedAt().After(time.Now().AddDate(0, -3, 0)) {
						_, _, _, err := client.Repositories.GetContents(ctx, repo.GetOwner().GetLogin(), repo.GetName(), ".gitpod.yml", &github.RepositoryContentGetOptions{})
						if err == nil {
							fmt.Printf("Org: %s, Repo: %s\n", org.GetLogin(), repo.GetName())
						}
					}
				}

				if resp.NextPage == 0 {
					break
				}
				repoOpt.Page = resp.NextPage
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
}

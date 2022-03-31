package service

import (
	"fmt"
	"net/url"

	"getsturdy.com/api/pkg/github"
	"getsturdy.com/api/pkg/users"
)

type Referer interface {
	URL() *url.URL
}

type gitHubPullRequestReferer struct {
	pr *github.PullRequest
}

func GitHubPullRequestReferer(pr *github.PullRequest) *gitHubPullRequestReferer {
	return &gitHubPullRequestReferer{pr: pr}
}

func (ghpr *gitHubPullRequestReferer) URL() *url.URL {
	return &url.URL{
		Scheme: "referer",
		Host:   "github",
		Path:   fmt.Sprintf("%d/prs/%s", ghpr.pr.GitHubRepositoryID, ghpr.pr.ID),
	}
}

type userReferer struct {
	u *users.User
}

func UserReferer(u *users.User) *userReferer {
	return &userReferer{u: u}
}

func (ur *userReferer) URL() *url.URL {
	return &url.URL{
		Scheme: "referer",
		Host:   "users",
		Path:   ur.u.ID.String(),
	}
}

package service

import (
	"fmt"
	"net/url"

	"getsturdy.com/api/pkg/github"
	"getsturdy.com/api/pkg/users"
)

type Referer interface {
	URL() string
}

type gitHubPullRequestReferer struct {
	pr *github.PullRequest
}

func GitHubPullRequestReferer(pr *github.PullRequest) *gitHubPullRequestReferer {
	return &gitHubPullRequestReferer{pr: pr}
}

func (ghpr *gitHubPullRequestReferer) URL() string {
	u := url.URL{
		Scheme: "referer",
		Host:   "github",
		Path:   fmt.Sprintf("%d/prs/%s", ghpr.pr.GitHubRepositoryID, ghpr.pr.ID),
	}
	return u.String()
}

type userReferer struct {
	u *users.User
}

func UserReferer(u *users.User) *userReferer {
	return &userReferer{u: u}
}

func (ur *userReferer) URL() string {
	u := url.URL{
		Scheme: "referer",
		Host:   "users",
		Path:   ur.u.ID.String(),
	}
	return u.String()
}

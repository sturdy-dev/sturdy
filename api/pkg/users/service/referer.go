package service

import (
	"fmt"
	"net/url"

	"getsturdy.com/api/pkg/users"
)

type Referer interface {
	URL() string
}

type gitHubPullRequestReferer struct {
	repoID int64
	prID   int64
}

func GitHubPullRequestReferer(repoID, prID int64) *gitHubPullRequestReferer {
	return &gitHubPullRequestReferer{repoID: repoID, prID: prID}
}

func (ghpr *gitHubPullRequestReferer) URL() string {
	u := url.URL{
		Scheme: "referer",
		Host:   "github",
		Path:   fmt.Sprintf("%d/prs/%d", ghpr.repoID, ghpr.prID),
	}
	return u.String()
}

type userReferer struct {
	userID users.ID
}

func UserReferer(userID users.ID) *userReferer {
	return &userReferer{userID: userID}
}

func (ur *userReferer) URL() string {
	u := url.URL{
		Scheme: "referer",
		Host:   "users",
		Path:   ur.userID.String(),
	}
	return u.String()
}

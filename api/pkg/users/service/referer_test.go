package service_test

import (
	"testing"

	"getsturdy.com/api/pkg/github"
	"getsturdy.com/api/pkg/users"
	service_users "getsturdy.com/api/pkg/users/service"
	"github.com/stretchr/testify/assert"
)

func Test_UserReferer(t *testing.T) {
	u := &users.User{ID: "testid"}
	assert.Equal(t, "referer://users/testid", service_users.UserReferer(u).URL())
}

func Test_GitHubPullRequestReferer(t *testing.T) {
	pr := &github.PullRequest{ID: "testid", GitHubRepositoryID: 1234}
	assert.Equal(t, "referer://github/1234/prs/testid", service_users.GitHubPullRequestReferer(pr).URL())
}

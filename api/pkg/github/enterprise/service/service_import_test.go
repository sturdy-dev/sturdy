package service_test

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
	"go.uber.org/zap"

	"getsturdy.com/api/pkg/auth"
	service_codebase "getsturdy.com/api/pkg/codebases/service"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/github"
	"getsturdy.com/api/pkg/github/api"
	service_github "getsturdy.com/api/pkg/github/enterprise/service"
	service_user "getsturdy.com/api/pkg/users/service"
	"getsturdy.com/api/pkg/workspaces"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/provider"
)

func TestService_ImportPullRequest(t *testing.T) {
	if os.Getenv("E2E_TEST") == "" {
		t.SkipNow()
	}

	type deps struct {
		dig.In

		GitHubService    *service_github.Service
		UserService      service_user.Service
		CodebaseService  *service_codebase.Service
		WorkspaceService service_workspace.Service

		RepoProvider provider.RepoProvider
		Logger       *zap.Logger
	}

	var d deps
	if !assert.NoError(t, di.Init(&d, module)) {
		t.FailNow()
	}

	ctx := context.Background()
	usr, err := d.UserService.CreateWithPassword(ctx, "hello", "foobar", "test+"+uuid.NewString()+"@getsturdy.com")

	authenticatedContext := auth.NewContext(ctx, &auth.Subject{ID: usr.ID.String(), Type: auth.SubjectUser})

	cb, err := d.CodebaseService.Create(authenticatedContext, "hello", nil)
	assert.NoError(t, err)

	// Create GitHub remote
	fakeGitHubBarePath := d.RepoProvider.ViewPath(cb.ID, "github")
	fakeGitHubBareRepo, err := vcs.CreateBareRepoWithRootCommit(fakeGitHubBarePath)
	_ = fakeGitHubBareRepo
	assert.NoError(t, err)

	trunkRepo, err := d.RepoProvider.TrunkRepo(cb.ID)
	assert.NoError(t, err)

	err = trunkRepo.AddNamedRemote("origin", fakeGitHubBarePath)
	assert.NoError(t, err)

	gitHubClonedPath := d.RepoProvider.ViewPath(cb.ID, "github-cloned")
	gitHubClonedRepo, err := vcs.CloneRepo(fakeGitHubBarePath, gitHubClonedPath)
	assert.NoError(t, err)

	err = gitHubClonedRepo.CreateNewBranchOnHEAD("new-pr")
	assert.NoError(t, err)
	err = gitHubClonedRepo.CheckoutBranchWithForce("new-pr")
	assert.NoError(t, err)
	err = ioutil.WriteFile(path.Join(gitHubClonedPath, "foobar.txt"), []byte("hello foobar"), 0644)
	assert.NoError(t, err)
	commitSha, err := gitHubClonedRepo.AddAndCommit("commit in pr")
	assert.NoError(t, err)
	err = gitHubClonedRepo.ForcePush(d.Logger, "new-pr")
	assert.NoError(t, err)
	err = fakeGitHubBareRepo.CreateRef("refs/pull/3/head", commitSha)
	assert.NoError(t, err)

	pr := &api.PullRequest{
		ID:     p[int64](2),
		Number: p[int](3),
		Title:  p[string]("hello"),
		Head:   &api.PullRequestBranch{SHA: &commitSha},
	}

	ghRepo := &github.Repository{
		ID:         uuid.NewString(),
		CodebaseID: cb.ID,
	}

	ghInstallation := &github.Installation{
		ID: uuid.NewString(),
	}

	err = d.GitHubService.ImportPullRequest(usr.ID, pr, ghRepo, ghInstallation, "testing-access-token")
	assert.NoError(t, err)

	// get workspace
	workspaceList, err := d.WorkspaceService.ListByCodebaseID(ctx, cb.ID, false)
	assert.NoError(t, err)

	var foundWs *workspaces.Workspace
	for _, ws := range workspaceList {
		if ws.DraftDescription == "<p>hello</p>" {
			foundWs = ws
		}
	}

	assert.NotNil(t, foundWs)
}

func p[T any](i T) *T {
	return &i
}

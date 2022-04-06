package service_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"testing"
	"time"

	gh "github.com/google/go-github/v39/github"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	service_codebase "getsturdy.com/api/pkg/codebases/service"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/github"
	"getsturdy.com/api/pkg/github/api"
	service_github "getsturdy.com/api/pkg/github/enterprise/service"
	service_github_webhooks "getsturdy.com/api/pkg/github/enterprise/webhooks"
	workers_github "getsturdy.com/api/pkg/github/enterprise/workers"
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

		GitHubService        *service_github.Service
		GitHubWebhookService *service_github_webhooks.Service
		GithubClonerQueue    *workers_github.ClonerQueue
		GithubImporterQueue  workers_github.ImporterQueue

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

	ctx, cancel := context.WithCancel(context.Background())
	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(func() error {
		if err := d.GithubClonerQueue.Start(ctx); err != nil {
			return fmt.Errorf("failed to start github cloner queue: %w", err)
		}
		return nil
	})

	wg.Go(func() error {
		if err := d.GithubImporterQueue.Start(ctx); err != nil {
			return fmt.Errorf("failed to start github importer queue: %w", err)
		}
		return nil
	})

	// cleanup workers
	defer func() {
		cancel()
		wg.Wait()
	}()

	rand.Seed(time.Now().UnixMilli())

	usr, err := d.UserService.CreateWithPassword(ctx, "hello", "foobar", "test+"+uuid.NewString()+"@getsturdy.com")

	// Create GitHub remote
	fakeGitHubBarePath := d.RepoProvider.ViewPath("not-a-codebase", "github-"+uuid.NewString())
	fakeGitHubBareRepo, err := vcs.CreateBareRepoWithRootCommit(fakeGitHubBarePath)
	assert.NoError(t, err)

	installation := &github.Installation{
		ID:             uuid.NewString(),
		InstallationID: rand.Int63n(600_000),
	}

	apiInstallation := &api.Installation{
		ID: &installation.InstallationID,
	}

	ghRepo := &gh.Repository{
		ID:       p[int64](rand.Int63n(800_00)),
		CloneURL: p[string](fakeGitHubBarePath),
	}

	sender := &gh.User{Email: p[string]("foobar")}

	FakeGitHubRepositoriesClient.repos[*ghRepo.ID] = *ghRepo

	cb, err := d.GitHubService.CreateNonReadyCodebaseAndClone(ctx, ghRepo, installation.InstallationID, sender, &usr.ID, nil)
	assert.NoError(t, err)

	repo := &github.Repository{
		ID:                 uuid.NewString(),
		CodebaseID:         cb.ID,
		InstallationID:     installation.InstallationID,
		GitHubRepositoryID: rand.Int63n(800_00),
	}

	apiRepo := api.ConvertRepository(ghRepo)

	gitHubClonedPath := d.RepoProvider.ViewPath(cb.ID, "github-cloned")
	gitHubClonedRepo, err := vcs.CloneRepo(fakeGitHubBarePath, gitHubClonedPath)
	assert.NoError(t, err)

	createFakePr := func(id int64, number int, title string) (pr *api.PullRequest, branchName string) {
		err = gitHubClonedRepo.CheckoutBranchWithForce("sturdytrunk")
		assert.NoError(t, err)
		branchName = fmt.Sprintf("pr-%s", uuid.NewString())
		err = gitHubClonedRepo.CreateNewBranchOnHEAD(branchName)
		assert.NoError(t, err)
		err = gitHubClonedRepo.CheckoutBranchWithForce(branchName)
		assert.NoError(t, err)
		err = ioutil.WriteFile(path.Join(gitHubClonedPath, "foobar.txt"), []byte("hello foobar "+branchName), 0644)
		assert.NoError(t, err)
		commitSha, err := gitHubClonedRepo.AddAndCommit("commit in pr")
		assert.NoError(t, err)
		err = gitHubClonedRepo.ForcePush(d.Logger, branchName)
		assert.NoError(t, err)
		err = fakeGitHubBareRepo.CreateRef(fmt.Sprintf("refs/pull/%d/head", number), commitSha)
		assert.NoError(t, err)
		pr = &api.PullRequest{
			ID:     &id,
			Number: &number,
			Title:  &title,
			Head:   &api.PullRequestBranch{SHA: &commitSha},
			State:  p[string](string(github.PullRequestStateOpen)),
		}
		return
	}

	getWorkspace := func(matchDescription string) *workspaces.Workspace {
		// get workspace
		workspaceList, err := d.WorkspaceService.ListByCodebaseID(ctx, cb.ID, false)
		assert.NoError(t, err)

		for _, ws := range workspaceList {
			log.Printf("ws: %+v", ws)
			if ws.DraftDescription == matchDescription {
				return ws
			}
		}
		return nil
	}

	t.Run("ImportPullRequest", func(t *testing.T) {
		pr, _ := createFakePr(5000, 1, "hello world")

		err = d.GitHubService.ImportPullRequest(usr.ID, pr, repo, installation, "testing-access-token")
		assert.NoError(t, err)

		ws := getWorkspace("<p>hello world</p>")
		if assert.NotNil(t, ws) {
			diffs, _, err := d.WorkspaceService.Diffs(ctx, ws.ID)
			assert.NoError(t, err)
			if assert.Len(t, diffs, 1) {
				assert.Equal(t, diffs[0].NewName, "foobar.txt")
			}
		}
	})

	t.Run("WebhookImport", func(t *testing.T) {
		pr, prBranchName := createFakePr(6000, 2, "hello webhook")

		err := d.GitHubWebhookService.HandlePullRequestEvent(ctx, &service_github_webhooks.PullRequestEvent{
			PullRequest:  pr,
			Repo:         apiRepo,
			Installation: apiInstallation,
		})
		assert.NoError(t, err)

		ws := getWorkspace("<p>hello webhook</p>")
		if assert.NotNil(t, ws) {
			diffs, _, err := d.WorkspaceService.Diffs(ctx, ws.ID)
			assert.NoError(t, err)
			if assert.Len(t, diffs, 1) {
				assert.Equal(t, diffs[0].NewName, "foobar.txt")
			}
		}

		afterFirstPushWorkspaceList, err := d.WorkspaceService.ListByCodebaseID(ctx, cb.ID, false)
		assert.NoError(t, err)

		// make changes to the pr
		err = gitHubClonedRepo.CheckoutBranchWithForce(prBranchName)
		assert.NoError(t, err)
		err = ioutil.WriteFile(path.Join(gitHubClonedPath, "foobar-other.txt"), []byte("hello foobar updated"), 0644)
		assert.NoError(t, err)
		commitSha, err := gitHubClonedRepo.AddAndCommit("commit 2 in pr")
		assert.NoError(t, err)
		err = gitHubClonedRepo.ForcePush(d.Logger, prBranchName)
		assert.NoError(t, err)
		err = fakeGitHubBareRepo.CreateRef(fmt.Sprintf("refs/pull/%d/head", pr.GetNumber()), commitSha)
		assert.NoError(t, err)
		pr.Head.SHA = &commitSha

		err = d.GitHubWebhookService.HandlePullRequestEvent(ctx, &service_github_webhooks.PullRequestEvent{
			PullRequest:  pr,
			Repo:         apiRepo,
			Installation: apiInstallation,
		})
		assert.NoError(t, err)

		afterSecondPushWorkspaceList, err := d.WorkspaceService.ListByCodebaseID(ctx, cb.ID, false)
		assert.NoError(t, err)

		if assert.NotNil(t, ws) {
			diffs, _, err := d.WorkspaceService.Diffs(ctx, ws.ID)
			assert.NoError(t, err)
			if assert.Len(t, diffs, 2) {
				assert.Equal(t, diffs[0].NewName, "foobar-other.txt")
				assert.Equal(t, diffs[1].NewName, "foobar.txt")
			}
		}

		// no new workspaces where created
		assert.Equal(t, len(afterFirstPushWorkspaceList), len(afterSecondPushWorkspaceList))

		// updated description and title on github
		pr.Body = p[string]("hello **body**")
		pr.Title = p[string]("this is a title")

		err = d.GitHubWebhookService.HandlePullRequestEvent(ctx, &service_github_webhooks.PullRequestEvent{
			PullRequest:  pr,
			Repo:         apiRepo,
			Installation: apiInstallation,
		})
		assert.NoError(t, err)

		updatedWs, err := d.WorkspaceService.GetByID(ctx, ws.ID)
		assert.NoError(t, err)

		assert.Equal(t, "<p>this is a title</p><p>hello <strong>body</strong></p>\n", updatedWs.DraftDescription)

		afterThirdPushWorkspaceList, err := d.WorkspaceService.ListByCodebaseID(ctx, cb.ID, false)
		assert.NoError(t, err)

		// no new workspaces where created
		assert.Equal(t, len(afterFirstPushWorkspaceList), len(afterThirdPushWorkspaceList))
	})
}

func p[T any](i T) *T {
	return &i
}

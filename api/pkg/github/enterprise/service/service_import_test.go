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
	"github.com/graph-gophers/graphql-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"getsturdy.com/api/pkg/auth"
	service_change "getsturdy.com/api/pkg/changes/service"
	service_codebase "getsturdy.com/api/pkg/codebases/service"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/github"
	"getsturdy.com/api/pkg/github/api"
	service_github "getsturdy.com/api/pkg/github/enterprise/service"
	service_github_webhooks "getsturdy.com/api/pkg/github/enterprise/webhooks"
	"getsturdy.com/api/pkg/graphql/resolvers"
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
		GithubClonerQueue    *service_github.ClonerQueue
		GithubImporterQueue  *service_github.ImporterQueue

		UserService      service_user.Service
		CodebaseService  *service_codebase.Service
		WorkspaceService service_workspace.Service
		ChangeService    *service_change.Service

		WorkspaceRootResolver resolvers.WorkspaceRootResolver

		RepoProvider provider.RepoProvider
		Logger       *zap.Logger
	}

	var d deps
	if !assert.NoError(t, di.Init(testModule).To(&d)) {
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
	assert.NoError(t, err)

	authenticatedCtx := auth.NewContext(ctx, &auth.Subject{Type: auth.SubjectUser, ID: string(usr.ID)})

	// Create GitHub remote
	fakeGitHubBarePath := d.RepoProvider.ViewPath("not-a-codebase", "github-"+uuid.NewString())
	fakeGitHubBareRepo, err := vcs.CreateBareRepoWithRootCommit(fakeGitHubBarePath)
	assert.NoError(t, err)
	err = fakeGitHubBareRepo.CreateNewBranchOnHEAD("master")
	assert.NoError(t, err)
	err = fakeGitHubBareRepo.SetDefaultBranch("master")
	assert.NoError(t, err)

	installation := &github.Installation{
		ID:             uuid.NewString(),
		InstallationID: rand.Int63n(600_000),
	}

	apiInstallation := &api.Installation{
		ID: &installation.InstallationID,
	}

	ghRepo := &gh.Repository{
		ID:            p[int64](rand.Int63n(800_00)),
		CloneURL:      p[string](fakeGitHubBarePath),
		DefaultBranch: p[string]("master"),
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
		TrackedBranch:      "master",
	}

	_, err = d.CodebaseService.AddUser(ctx, cb.ID, usr, usr.ID)
	assert.NoError(t, err)

	apiRepo := api.ConvertRepository(ghRepo)

	gitHubClonedPath := d.RepoProvider.ViewPath(cb.ID, "github-cloned")
	gitHubClonedRepo, err := vcs.CloneRepo(fakeGitHubBarePath, gitHubClonedPath)
	assert.NoError(t, err)

	createFakePr := func(id int64, number int, fileName, title string) (pr *api.PullRequest, branchName string) {
		err = gitHubClonedRepo.CheckoutBranchWithForce("master")
		assert.NoError(t, err)
		branchName = fmt.Sprintf("pr-%s", uuid.NewString())
		err = gitHubClonedRepo.CreateNewBranchOnHEAD(branchName)
		assert.NoError(t, err)
		err = gitHubClonedRepo.CheckoutBranchWithForce(branchName)
		assert.NoError(t, err)
		err = ioutil.WriteFile(path.Join(gitHubClonedPath, fileName), []byte("hello foobar "+fileName), 0644)
		assert.NoError(t, err)
		commitSha, err := gitHubClonedRepo.AddAndCommit("commit in pr")
		assert.NoError(t, err)
		err = gitHubClonedRepo.ForcePush(d.Logger, branchName)
		assert.NoError(t, err)
		err = fakeGitHubBareRepo.CreateRef(fmt.Sprintf("refs/pull/%d/head", number), commitSha)
		assert.NoError(t, err)

		masterCommitID, err := gitHubClonedRepo.BranchCommitID("master")
		assert.NoError(t, err)

		pr = &api.PullRequest{
			ID:     &id,
			Number: &number,
			Title:  &title,
			Head:   &api.PullRequestBranch{SHA: &commitSha},
			Base:   &api.PullRequestBranch{SHA: &masterCommitID},
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
		pr, _ := createFakePr(5000, 1, "foobar.txt", "hello world")

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
		pr, prBranchName := createFakePr(6000, 2, "foobar.txt", "hello webhook")

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

		// pr should be updatable
		{
			workspaceResolver, err := d.WorkspaceRootResolver.Workspace(authenticatedCtx, resolvers.WorkspaceArgs{ID: graphql.ID(ws.ID)})
			assert.NoError(t, err)
			gitHubPullRequestResolver, err := workspaceResolver.GitHubPullRequest(authenticatedCtx)
			assert.NoError(t, err)
			assert.Equal(t, int32(pr.GetNumber()), int32(gitHubPullRequestResolver.PullRequestNumber()))
			assert.True(t, gitHubPullRequestResolver.CanUpdate())
		}
	})

	t.Run("ImportedParent", func(t *testing.T) {
		// create two prs
		firstPR, firstBranchName := createFakePr(7000, 3, "first.txt", "hello first")
		secondPR, _ := createFakePr(7001, 4, "second.txt", "hello second")

		// import the first pr
		err = d.GitHubWebhookService.HandlePullRequestEvent(ctx, &service_github_webhooks.PullRequestEvent{
			PullRequest:  firstPR,
			Repo:         apiRepo,
			Installation: apiInstallation,
		})
		assert.NoError(t, err)

		ws := getWorkspace("<p>hello first</p>")
		assert.NotNil(t, ws)

		// pr should be updatable
		{
			workspaceResolver, err := d.WorkspaceRootResolver.Workspace(authenticatedCtx, resolvers.WorkspaceArgs{ID: graphql.ID(ws.ID)})
			assert.NoError(t, err)
			gitHubPullRequestResolver, err := workspaceResolver.GitHubPullRequest(authenticatedCtx)
			assert.NoError(t, err)
			assert.Equal(t, int32(firstPR.GetNumber()), int32(gitHubPullRequestResolver.PullRequestNumber()))
			assert.True(t, gitHubPullRequestResolver.CanUpdate())
		}

		// merge first pr
		mergeCommitSha, err := fakeGitHubBareRepo.MergeBranchInto(firstBranchName, "master")
		assert.NoError(t, err)

		firstPR.State = p[string]("closed")
		firstPR.Merged = p[bool](true)
		firstPR.MergeCommitSHA = &mergeCommitSha

		err = d.GitHubWebhookService.HandlePullRequestEvent(ctx, &service_github_webhooks.PullRequestEvent{
			PullRequest:  firstPR,
			Repo:         apiRepo,
			Installation: apiInstallation,
		})
		assert.NoError(t, err)

		// check merged result
		cb, err = d.CodebaseService.GetByID(ctx, cb.ID) // reload
		assert.NoError(t, err)
		headChange, err := d.ChangeService.HeadChange(ctx, cb)
		assert.NoError(t, err)
		assert.Equal(t, "hello first", *headChange.Title)

		// open second pr
		err = d.GitHubWebhookService.HandlePullRequestEvent(ctx, &service_github_webhooks.PullRequestEvent{
			PullRequest:  secondPR,
			Repo:         apiRepo,
			Installation: apiInstallation,
		})
		assert.NoError(t, err)

		secondWorkspace := getWorkspace("<p>hello second</p>")
		assert.NotNil(t, secondWorkspace)

		diffs, _, err := d.WorkspaceService.Diffs(ctx, secondWorkspace.ID)
		assert.NotNil(t, secondWorkspace)
		t.Logf("diffs: %+v", diffs)
		if assert.Len(t, diffs, 1) {
			assert.Equal(t, "second.txt", diffs[0].NewName)
		}
	})

	t.Run("ImportFromFork", func(t *testing.T) {
		pr, _ := createFakePr(8000, 8, "foobar.txt", "forked PR")

		// some random user
		pr.Head.User = &api.User{
			ID:    p[int64](rand.Int63n(1_000_00)),
			Login: p[string](uuid.NewString()),
			Email: p[string](uuid.NewString() + "@testing.getsturdy.com"),
			Name:  p[string](uuid.NewString()),
		}

		err := d.GitHubWebhookService.HandlePullRequestEvent(ctx, &service_github_webhooks.PullRequestEvent{
			PullRequest:  pr,
			Repo:         apiRepo,
			Installation: apiInstallation,
		})
		assert.NoError(t, err)

		ws := getWorkspace("<p>forked PR</p>")
		if assert.NotNil(t, ws) {
			diffs, _, err := d.WorkspaceService.Diffs(ctx, ws.ID)
			assert.NoError(t, err)
			if assert.Len(t, diffs, 1) {
				assert.Equal(t, diffs[0].NewName, "foobar.txt")
			}
		}

		afterFirstPushWorkspaceList, err := d.WorkspaceService.ListByCodebaseID(ctx, cb.ID, false)
		assert.NoError(t, err)

		// send event again
		err = d.GitHubWebhookService.HandlePullRequestEvent(ctx, &service_github_webhooks.PullRequestEvent{
			PullRequest:  pr,
			Repo:         apiRepo,
			Installation: apiInstallation,
		})
		assert.NoError(t, err)

		afterSecondPushWorkspaceList, err := d.WorkspaceService.ListByCodebaseID(ctx, cb.ID, false)
		assert.NoError(t, err)
		assert.Len(t, afterSecondPushWorkspaceList, len(afterFirstPushWorkspaceList))

		// pr is from a fork, and should not be updatable
		{
			workspaceResolver, err := d.WorkspaceRootResolver.Workspace(authenticatedCtx, resolvers.WorkspaceArgs{ID: graphql.ID(ws.ID)})
			assert.NoError(t, err)
			gitHubPullRequestResolver, err := workspaceResolver.GitHubPullRequest(authenticatedCtx)
			assert.NoError(t, err)
			assert.Equal(t, int32(pr.GetNumber()), int32(gitHubPullRequestResolver.PullRequestNumber()))
			assert.False(t, gitHubPullRequestResolver.CanUpdate())
		}
	})
}

func p[T any](i T) *T {
	return &i
}

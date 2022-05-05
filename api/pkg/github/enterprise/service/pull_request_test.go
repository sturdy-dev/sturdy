//nolint:bodyclose
package service_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	gh "github.com/google/go-github/v39/github"
	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	activity_sender "getsturdy.com/api/pkg/activity/sender"
	analytics_service "getsturdy.com/api/pkg/analytics/service"
	module_api "getsturdy.com/api/pkg/api/module"
	"getsturdy.com/api/pkg/auth"
	workers_ci "getsturdy.com/api/pkg/ci/workers"
	"getsturdy.com/api/pkg/codebases"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	service_comments "getsturdy.com/api/pkg/comments/service"
	"getsturdy.com/api/pkg/configuration"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/events"
	"getsturdy.com/api/pkg/github"
	"getsturdy.com/api/pkg/github/enterprise/client"
	"getsturdy.com/api/pkg/github/enterprise/config"
	db_github "getsturdy.com/api/pkg/github/enterprise/db"
	graphql_github "getsturdy.com/api/pkg/github/enterprise/graphql"
	graphql_github_pr "getsturdy.com/api/pkg/github/enterprise/graphql/pr"
	"getsturdy.com/api/pkg/github/enterprise/routes"
	service_github "getsturdy.com/api/pkg/github/enterprise/service"
	service_github_webhooks "getsturdy.com/api/pkg/github/enterprise/webhooks"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	module_queue "getsturdy.com/api/pkg/queue/module"
	db_review "getsturdy.com/api/pkg/review/db"
	service_statuses "getsturdy.com/api/pkg/statuses/service"
	service_sync "getsturdy.com/api/pkg/sync/service"
	"getsturdy.com/api/pkg/users"
	db_user "getsturdy.com/api/pkg/users/db"
	"getsturdy.com/api/pkg/view"
	db_view "getsturdy.com/api/pkg/view/db"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"
	"getsturdy.com/api/vcs/provider"
)

func testModule(t *testing.T) di.Module {
	return func(c *di.Container) {
		c.Import(module_api.Module)

		c.Register(service_github_webhooks.New)
		c.Register(service_github_webhooks.NewWebhooksQueue)

		c.ImportWithForce(service_github.Module)
		c.ImportWithForce(graphql_github.Module)
		c.ImportWithForce(graphql_github_pr.Module)
		c.ImportWithForce(configuration.TestModule)
		c.ImportWithForce(module_queue.TestModule(t))

		c.RegisterWithForce(func() client.InstallationClientProvider { return clientProvider })
		c.RegisterWithForce(func() client.PersonalClientProvider { return personalClientProvider })
		c.RegisterWithForce(func() client.AppClientProvider { return appsClientProvider })
		c.RegisterWithForce(func() *config.GitHubAppConfig { return &config.GitHubAppConfig{} })
	}
}

func TestPRHighLevel(t *testing.T) {
	if os.Getenv("E2E_TEST") == "" {
		t.SkipNow()
	}

	type deps struct {
		di.In

		ActivitySender                activity_sender.ActivitySender
		AnalyticsClient               *analytics_service.Service
		BuildQueue                    *workers_ci.BuildQueue
		CodebaseRootResolver          resolvers.CodebaseRootResolver
		CodebaseRepo                  db_codebases.CodebaseRepository
		CodebaseUserRepo              db_codebases.CodebaseUserRepository
		CommentsRootResolver          resolvers.CommentRootResolver
		CommentsService               *service_comments.Service
		EventsSender                  events.EventSender
		ExecutorProvider              executor.Provider
		GitHubInstallationRepo        db_github.GitHubInstallationRepository
		GitHubPRRepo                  db_github.GitHubPRRepository
		GitHubPullRequestRootResolver resolvers.GitHubPullRequestRootResolver
		GitHubRepositoryRepo          db_github.GitHubRepositoryRepository
		GitHubService                 *service_github.Service
		GitHubWebhookService          *service_github_webhooks.Service
		GitHubUserRepo                db_github.GitHubUserRepository
		Logger                        *zap.Logger
		RepoProvider                  provider.RepoProvider
		ReviewsRepo                   db_review.ReviewRepository
		StatusesService               *service_statuses.Service
		SyncService                   *service_sync.Service
		UserRepo                      db_user.Repository
		ViewRepo                      db_view.Repository
		ViewRootResolver              resolvers.ViewRootResolver
		WorkspaceRepo                 db_workspaces.Repository
		WorkspaceRootResolver         resolvers.WorkspaceRootResolver
		WorkspaceService              *service_workspace.Service
		WebhooksQueue                 *service_github_webhooks.Queue
	}

	var d deps
	if !assert.NoError(t, di.Init(testModule(t)).To(&d)) {
		t.FailNow()
	}

	commentsResolver := d.CommentsRootResolver
	viewResolver := d.ViewRootResolver
	gitHubRepositoryRepo := d.GitHubRepositoryRepo
	gitHubUserRepo := d.GitHubUserRepo
	gitHubInstallationRepo := d.GitHubInstallationRepo
	codebaseRepo := d.CodebaseRepo
	gitHubPRRepo := d.GitHubPRRepo
	workspaceResolver := d.WorkspaceRootResolver
	prResolver := d.GitHubPullRequestRootResolver
	workspaceService := d.WorkspaceService
	repoProvider := d.RepoProvider
	codebaseUserRepo := d.CodebaseUserRepo
	workspaceRepo := d.WorkspaceRepo
	viewRepo := d.ViewRepo

	webhookRoute := routes.Webhook(d.Logger, d.WebhooksQueue)

	go func() {
		assert.NoError(t, d.WebhooksQueue.Start(context.TODO()))
	}()
	go func() {
		assert.NoError(t, d.BuildQueue.Start(context.TODO()))
	}()

	testCases := []struct {
		name                       string
		gitHubRebase               bool
		expectedHunkID             string
		changeFiles                map[string]string
		withCommitsAlreadyOnGitHub bool
	}{
		{
			name:         "rebase",
			gitHubRebase: true,
			changeFiles: map[string]string{
				"a.txt": "foo\nbar\nbaz2\n",
				"b.txt": "bbb\nbbb\nbbb\nBBB\n",
			},
			expectedHunkID:             "aaaa",
			withCommitsAlreadyOnGitHub: true,
		},
		{
			name:         "rebase-CRLF",
			gitHubRebase: true,
			changeFiles: map[string]string{
				"a.txt": "foo\r\nbar\r\nbaz2\r\n",
				"b.txt": "bbb\r\nbbb\r\nbbb\r\nBBB\r\n",
			},
			expectedHunkID:             "aaaa",
			withCommitsAlreadyOnGitHub: true,
		},
		{
			name:         "rebase-b-txt-remove-newline",
			gitHubRebase: true,
			changeFiles: map[string]string{
				"a.txt": "foo\nbar\nbaz2\n",
				"b.txt": "bbb\nbbb\nbbb",
			},
			expectedHunkID:             "bbbb",
			withCommitsAlreadyOnGitHub: true,
		},

		{
			name:         "merge",
			gitHubRebase: false,
			changeFiles: map[string]string{
				"a.txt": "foo\nbar\nbaz2\n",
				"b.txt": "bbb\nbbb\nbbb\nBBB\n",
			},
			expectedHunkID:             "aaaa",
			withCommitsAlreadyOnGitHub: true,
		},
		{
			name:         "merge-CRLF",
			gitHubRebase: false,
			changeFiles: map[string]string{
				"a.txt": "foo\r\nbar\r\nbaz2\r\n",
				"b.txt": "bbb\r\nbbb\r\nbbb\r\nBBB\r\n",
			},
			expectedHunkID:             "aaaa",
			withCommitsAlreadyOnGitHub: true,
		},
		{
			name:         "merge-CRLF-remove-trailing",
			gitHubRebase: false,
			changeFiles: map[string]string{
				"a.txt": "foo\r\nbar\r\nbaz2\r\n",
				"b.txt": "bbb\r\nbbb\r\nbbb",
			},
			expectedHunkID:             "aaaa",
			withCommitsAlreadyOnGitHub: true,
		},

		{
			name:         "merge-github-empty-clone",
			gitHubRebase: false,
			changeFiles: map[string]string{
				"a.txt": "foo\nbar\nbaz2\n",
				"b.txt": "bbb\nbbb\nbbb\nBBB\n",
			},
			expectedHunkID:             "aaaa",
			withCommitsAlreadyOnGitHub: false,
		},
		{
			name:         "rebase-github-empty-clone",
			gitHubRebase: true,
			changeFiles: map[string]string{
				"a.txt": "foo\nbar\nbaz2\n",
				"b.txt": "bbb\nbbb\nbbb\nBBB\n",
			},
			expectedHunkID:             "aaaa",
			withCommitsAlreadyOnGitHub: false,
		},
	}

	rand.Seed(time.Now().UnixNano())

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			userID := users.ID(uuid.NewString())
			viewID := uuid.NewString()
			codebaseID := codebases.ID(uuid.NewString())
			codebaseUserID := uuid.NewString()
			sturdyRepositoryID := uuid.NewString()
			gitHubRepositoryID := rand.Int63n(500_000_000)
			gitHubInstallationID := rand.Int63n(500_000_000)
			ctx := auth.NewContext(context.Background(), &auth.Subject{Type: auth.SubjectUser, ID: userID.String()})

			gitHubRepoOwner := uuid.NewString()
			gitHubRepoName := uuid.NewString()

			cu := &codebases.CodebaseUser{
				ID:         codebaseUserID,
				UserID:     userID,
				CodebaseID: codebaseID,
			}
			expT := time.Now().Add(20 * time.Minute)
			ghr := &github.Repository{
				ID:                               sturdyRepositoryID,
				GitHubRepositoryID:               gitHubRepositoryID,
				InstallationID:                   gitHubInstallationID,
				Name:                             gitHubRepoName,
				TrackedBranch:                    "master",
				CodebaseID:                       codebaseID,
				GitHubSourceOfTruth:              true,
				IntegrationEnabled:               true,
				InstallationAccessToken:          p("token"),
				InstallationAccessTokenExpiresAt: &expT,
			}
			ghu := &github.User{
				ID:          uuid.NewString(),
				UserID:      userID,
				Username:    uuid.NewString(),
				AccessToken: p("let's assume this is valid"),
			}

			in := &github.Installation{
				ID:             uuid.NewString(),
				InstallationID: gitHubInstallationID,
				Owner:          gitHubRepoOwner,
			}

			assert.NoError(t, d.UserRepo.Create(&users.User{ID: userID, Email: userID.String() + "@getsturdy.com", Name: "Test Testsson", Status: users.StatusActive}))
			assert.NoError(t, codebaseRepo.Create(codebases.Codebase{ID: codebaseID, ShortCodebaseID: codebases.ShortCodebaseID(codebaseID)}))
			assert.NoError(t, codebaseUserRepo.Create(*cu))
			assert.NoError(t, gitHubUserRepo.Create(*ghu))
			assert.NoError(t, gitHubInstallationRepo.Create(*in))
			assert.NoError(t, gitHubRepositoryRepo.Create(*ghr))

			// Create GitHub remote
			var err error
			fakeGitHubRemotePath := repoProvider.ViewPath(codebaseID, "github")
			var fakeGitHubBareRepo vcs.RepoWriter
			if tc.withCommitsAlreadyOnGitHub {
				fakeGitHubBareRepo, err = vcs.CreateBareRepoWithRootCommit(fakeGitHubRemotePath)
			} else {
				fakeGitHubBareRepo, err = vcs.CreateEmptyBareRepo(fakeGitHubRemotePath)
			}
			if !assert.NoError(t, err) {
				t.FailNow()
			}

			// Clone to the trunk
			pathBase := repoProvider.TrunkPath(codebaseID)
			t.Logf("base=%s", pathBase)

			_, err = vcs.CloneRepoBare(fakeGitHubRemotePath, pathBase)
			if !assert.NoError(t, err) {
				t.FailNow()
			}

			// Create workspace
			wsRes, err := workspaceResolver.CreateWorkspace(ctx, resolvers.CreateWorkspaceArgs{Input: resolvers.CreateWorkspaceInput{CodebaseID: graphql.ID(codebaseID)}})
			assert.NoError(t, err)
			workspaceID := string(wsRes.ID())

			vw := &view.View{
				ID:          viewID,
				CodebaseID:  codebaseID,
				UserID:      userID,
				WorkspaceID: workspaceID,
			}
			assert.NoError(t, viewRepo.Create(*vw))

			// Clone to the view
			viewApath := repoProvider.ViewPath(codebaseID, viewID)
			_, err = vcs.CloneRepo(pathBase, viewApath)
			if !assert.NoError(t, err) {
				t.FailNow()
			}

			// Open the workspace on the view
			_, err = viewResolver.OpenWorkspaceOnView(ctx, resolvers.OpenViewArgs{Input: resolvers.OpenWorkspaceOnViewInput{
				WorkspaceID: graphql.ID(workspaceID),
				ViewID:      graphql.ID(viewID),
			}})
			if !assert.NoError(t, err) {
				//nolint:errorlint
				t.Logf("err: %+v", err.(*gqlerrors.SturdyGraphqlError).OriginalError())
			}

			repo, err := repoProvider.ViewRepo(codebaseID, viewID)
			assert.NoError(t, err)
			headBranchName, err := repo.HeadBranch()
			assert.NoError(t, err)
			assert.Equal(t, workspaceID, headBranchName)

			// setup complete

			// make changes
			for name, content := range tc.changeFiles {
				assert.NoError(t, ioutil.WriteFile(path.Join(viewApath, name), []byte(content), 0666))
			}

			// Set workspace draft description
			workspaceIDgql := graphql.ID(workspaceID)
			_, err = workspaceResolver.UpdateWorkspace(ctx, resolvers.UpdateWorkspaceArgs{
				Input: resolvers.UpdateWorkspaceInput{
					ID:               workspaceIDgql,
					DraftDescription: p("<p><em>draft description</em></p>"),
				},
			})
			assert.NoError(t, err)

			// Add comments on the workspace
			viewIDgql := graphql.ID(viewID)
			for i := 0; i < 5; i++ {
				_, err = commentsResolver.CreateComment(ctx, resolvers.CreateCommentArgs{Input: resolvers.CreateCommentInput{
					Message:     fmt.Sprintf("commenting on a workspace i=%d", i),
					WorkspaceID: &workspaceIDgql,
					ViewID:      &viewIDgql,
					Path:        p("a.txt"),
					LineStart:   p[int32](1),
					LineEnd:     p[int32](1),
					LineIsNew:   p(true),
				}})
				assert.NoError(t, err)
				//nolint:errorlint
				if gerr, ok := err.(*gqlerrors.SturdyGraphqlError); ok {
					assert.NoError(t, gerr.OriginalError())
				}
			}

			// Get comments from workspace
			wsResolver, err := workspaceResolver.Workspace(ctx, resolvers.WorkspaceArgs{ID: workspaceIDgql})
			assert.NoError(t, err)
			workspaceComments, err := wsResolver.Comments()
			assert.NoError(t, err)
			assert.Len(t, workspaceComments, 5)

			// Get diffs
			diffs, _, err := workspaceService.Diffs(context.Background(), workspaceID)
			assert.NoError(t, err)
			t.Logf("diffs=%+v", diffs)

			var hunkIDs []string
			for _, diff := range diffs {
				for _, hunk := range diff.Hunks {
					hunkIDs = append(hunkIDs, hunk.ID)
				}
			}
			assert.NotEmpty(t, hunkIDs)

			// Create initial Pull request
			createdPullRequestResolver, err := prResolver.CreateOrUpdateGitHubPullRequest(ctx,
				resolvers.CreateOrUpdateGitHubPullRequestArgs{
					Input: resolvers.CreateOrUpdateGitHubPullRequestInput{
						WorkspaceID: graphql.ID(workspaceID),
					}},
			)
			if !assert.NoError(t, err, errors.Unwrap(err)) {
				//nolint:errorlint
				t.Logf("err=%+v", err.(*gqlerrors.SturdyGraphqlError).OriginalError())
			}
			if assert.NotNil(t, createdPullRequestResolver) {
				assert.True(t, createdPullRequestResolver.Open())
				assert.False(t, createdPullRequestResolver.Merged())
				if status, err := createdPullRequestResolver.State(); assert.NoError(t, err) {
					assert.Equal(t, resolvers.GitHubPullRequestStateOpen, status)
				}
			} else {
				t.FailNow()
			}

			// get githubs pull request ID (not the same as the pull request number)
			ghpr, err := gitHubPRRepo.Get(string(createdPullRequestResolver.ID()))
			assert.NoError(t, err)
			gitHubPullRequestID := ghpr.GitHubID

			// PR was closed
			prWebhookEvent(t, userID, webhookRoute, gh.PullRequestEvent{
				PullRequest: &gh.PullRequest{
					ID:    &gitHubPullRequestID,
					State: p("closed"),
				},
				Repo:         &gh.Repository{ID: &gitHubRepositoryID},
				Installation: &gh.Installation{ID: &gitHubInstallationID},
			})

			// Updated PR is closed
			gqlID := graphql.ID(workspaceID)
			updatedPR, err := prResolver.InternalGitHubPullRequestByWorkspaceID(ctx, resolvers.GitHubPullRequestArgs{WorkspaceID: &gqlID})
			assert.NoError(t, err)
			assert.False(t, updatedPR.Merged())
			assert.False(t, updatedPR.Open())
			if status, err := updatedPR.State(); assert.NoError(t, err) {
				assert.Equal(t, resolvers.GitHubPullRequestStateClosed, status)
			}

			// PR was reopened
			prWebhookEvent(t, userID, webhookRoute, gh.PullRequestEvent{
				PullRequest: &gh.PullRequest{
					ID:    &gitHubPullRequestID,
					State: p("open"),
				},
				Repo:         &gh.Repository{ID: &gitHubRepositoryID},
				Installation: &gh.Installation{ID: &gitHubInstallationID},
			})

			// Updated PR is opened
			gqlID = graphql.ID(workspaceID)
			updatedPR, err = prResolver.InternalGitHubPullRequestByWorkspaceID(ctx, resolvers.GitHubPullRequestArgs{WorkspaceID: &gqlID})
			assert.NoError(t, err)
			assert.True(t, updatedPR.Open())
			if status, err := updatedPR.State(); assert.NoError(t, err) {
				assert.Equal(t, resolvers.GitHubPullRequestStateOpen, status)
			}

			preMergeHeadSha, err := fakeGitHubBareRepo.BranchCommitID("master")
			assert.NoError(t, err)

			// Rebase or rebase the commit
			var mergeCommitSha string
			if tc.gitHubRebase {
				branchCommit, err := fakeGitHubBareRepo.BranchCommitID("sturdy-pr-" + workspaceID)
				assert.NoError(t, err)
				mergeCommitSha, _, _, err = fakeGitHubBareRepo.CherryPickOnto(branchCommit, preMergeHeadSha)
				assert.NoError(t, err)
				err = fakeGitHubBareRepo.MoveBranchToCommit("master", mergeCommitSha)
				assert.NoError(t, err)
			} else {
				mergeCommitSha, err = fakeGitHubBareRepo.MergeBranchInto("sturdy-pr-"+workspaceID, "master")
				assert.NoError(t, err)
			}

			// Merge PR
			prWebhookEvent(t, userID, webhookRoute, gh.PullRequestEvent{
				PullRequest: &gh.PullRequest{
					ID:             &gitHubPullRequestID,
					State:          p("closed"),
					Merged:         p(true),
					MergeCommitSHA: &mergeCommitSha,
					Base: &gh.PullRequestBranch{
						SHA: &preMergeHeadSha,
					},
				},
				Repo:         &gh.Repository{ID: &gitHubRepositoryID},
				Installation: &gh.Installation{ID: &gitHubInstallationID},
			})

			// Updated PR is merged
			updatedPR, err = prResolver.InternalGitHubPullRequestByWorkspaceID(ctx, resolvers.GitHubPullRequestArgs{WorkspaceID: &gqlID})
			assert.NoError(t, err)
			assert.False(t, updatedPR.Open())
			assert.True(t, updatedPR.Merged())
			if status, err := updatedPR.State(); assert.NoError(t, err) {
				assert.Equal(t, resolvers.GitHubPullRequestStateMerged, status)
			}

			// Post-merge push webhook event
			webhookRepoPush := gh.PushEvent{
				Ref:          p("refs/heads/master"),
				Repo:         &gh.PushEventRepository{ID: &gitHubRepositoryID},
				Installation: &gh.Installation{ID: &gitHubInstallationID},
			}
			requestWithParams(t, userID, webhookRoute, webhookRepoPush, nil, "push", []gin.Param{})

			// Workspace up to date state is reset after the push event
			ws, err := workspaceRepo.Get(workspaceID)
			assert.NoError(t, err)
			assert.Nil(t, ws.UpToDateWithTrunk)

			// The workspace should no longer have any comments
			wsResolver, err = workspaceResolver.Workspace(ctx, resolvers.WorkspaceArgs{ID: gqlID, AllowArchived: p(true)})
			assert.NoError(t, err)
			workspaceComments, err = wsResolver.Comments()
			assert.NoError(t, err)
			assert.Len(t, workspaceComments, 0)

			gqlCodebaseID := graphql.ID(codebaseID)
			codebaseResolver, err := d.CodebaseRootResolver.Codebase(ctx, resolvers.CodebaseArgs{ID: &gqlCodebaseID})
			assert.NoError(t, err)
			changeResolvers, err := codebaseResolver.Changes(ctx, &resolvers.CodebaseChangesArgs{Input: &resolvers.CodebaseChangesInput{Limit: p[int32](50)}})
			assert.NoError(t, err)

			var found bool
			for _, changeResolver := range changeResolvers {
				if changeResolver.Title() == "draft description" {
					found = true

					author, err := changeResolver.Author(ctx)
					assert.NoError(t, err)
					assert.Equal(t, "Test Testsson", author.Name())

					assert.Equal(t, "<p><em>draft description</em></p>", changeResolver.Description())

					// The new change should have comments
					changeComments, err := changeResolver.Comments()
					assert.NoError(t, err)
					assert.Len(t, changeComments, 5)
				}
			}
			assert.True(t, found)
		})
	}
}

func prWebhookEvent(t *testing.T, userID users.ID, webhookRoute gin.HandlerFunc, event gh.PullRequestEvent) {
	requestWithParams(t, userID, webhookRoute, event, nil, "pull_request", []gin.Param{})
}

func clientProvider(gitHubAppConfig *config.GitHubAppConfig, installationID int64) (tokenClient *client.GitHubClients, appsClient client.AppsClient, err error) {
	return &client.GitHubClients{
			Repositories: FakeGitHubRepositoriesClient,
			PullRequests: &fakeGitHubPullRequestClient{},
		},
		&fakeGitHubAppsClient{}, nil
}

func appsClientProvider(gitHubAppConfig *config.GitHubAppConfig) (client.AppsClient, error) {
	return &fakeGitHubAppsClient{}, nil
}

func personalClientProvider(token string) (*client.GitHubClients, error) {
	return &client.GitHubClients{
		Repositories: nil,
		PullRequests: &fakeGitHubPullRequestClient{},
	}, nil
}

type fakeGitHubPullRequestClient struct {
	prs []*gh.PullRequest
}

func (f *fakeGitHubPullRequestClient) List(ctx context.Context, owner string, repo string, opts *gh.PullRequestListOptions) ([]*gh.PullRequest, *gh.Response, error) {
	panic("implement me")
}

func (f *fakeGitHubPullRequestClient) Create(ctx context.Context, owner string, repo string, pull *gh.NewPullRequest) (*gh.PullRequest, *gh.Response, error) {
	rand.Seed(time.Now().UnixNano())
	id := int64(rand.Intn(10000))
	num := rand.Intn(10000)
	pr := gh.PullRequest{
		ID:     &id,
		Number: &num,
		State:  p("open"),
		Title:  pull.Title,
		Body:   pull.Body,
		Head:   &gh.PullRequestBranch{Ref: pull.Head},
		Base:   &gh.PullRequestBranch{Ref: pull.Base},
	}
	f.prs = append(f.prs, &pr)
	return &pr, nil, nil
}

func (f *fakeGitHubPullRequestClient) Get(ctx context.Context, owner string, repo string, number int) (*gh.PullRequest, *gh.Response, error) {
	panic("implement me")
}

func (f *fakeGitHubPullRequestClient) Edit(ctx context.Context, owner string, repo string, number int, pull *gh.PullRequest) (*gh.PullRequest, *gh.Response, error) {
	panic("implement me")
}

type fakeGitHubAppsClient struct{}

func (f *fakeGitHubAppsClient) CreateInstallationToken(ctx context.Context, id int64, opts *gh.InstallationTokenOptions) (*gh.InstallationToken, *gh.Response, error) {
	return &gh.InstallationToken{
		Token:        p("testingtoken"),
		ExpiresAt:    p(time.Now().Add(time.Hour * 3)),
		Permissions:  opts.Permissions,
		Repositories: nil,
	}, nil, nil
}

func (f *fakeGitHubAppsClient) GetInstallation(ctx context.Context, id int64) (*gh.Installation, *gh.Response, error) {
	return &gh.Installation{
		ID: &id,
	}, nil, nil
}

func (f *fakeGitHubAppsClient) Get(ctx context.Context, appSlug string) (*gh.App, *gh.Response, error) {
	panic("implement me")
}

var FakeGitHubRepositoriesClient = &fakeGitHubRepositoriesClient{
	repos: make(map[int64]gh.Repository),
}

type fakeGitHubRepositoriesClient struct {
	repos map[int64]gh.Repository
}

func (f fakeGitHubRepositoriesClient) Get(ctx context.Context, owner, repo string) (*gh.Repository, *gh.Response, error) {
	for _, r := range f.repos {
		if r.GetOwner().GetLogin() == owner && r.GetName() == repo {
			return &r, nil, nil
		}
	}
	return nil, nil, fmt.Errorf("fakeGitHubRepositoriesClient.Get: not found")
}

func (f fakeGitHubRepositoriesClient) GetByID(ctx context.Context, id int64) (*gh.Repository, *gh.Response, error) {
	if r, ok := f.repos[id]; ok {
		return &r, nil, nil
	}
	return nil, nil, fmt.Errorf("fakeGitHubRepositoriesClient.GetByID: not found")
}

func (f fakeGitHubRepositoriesClient) ListCollaborators(ctx context.Context, owner, repo string, opts *gh.ListCollaboratorsOptions) ([]*gh.User, *gh.Response, error) {
	return []*gh.User{}, &gh.Response{
		FirstPage: 1,
		LastPage:  1,
		NextPage:  0,
	}, nil
}

func requestWithParams(t *testing.T, userID users.ID, route func(*gin.Context), request, response any, reqType string, params []gin.Param) {
	res := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(res)
	c.Params = params

	data, err := json.Marshal(request)
	assert.NoError(t, err)

	c.Request, err = http.NewRequest("GET", "/", bytes.NewReader(data))
	c.Request = c.Request.WithContext(auth.NewContext(context.Background(), &auth.Subject{ID: userID.String(), Type: auth.SubjectUser}))
	assert.NoError(t, err)
	c.Request.Header.Set("X-Hub-Signature", "sha1=126f2c800419c60137ce748d7672e77b65cf16d6")
	c.Request.Header.Set("X-Github-Event", reqType)
	c.Request.Header.Set("Content-Type", "application/json")

	assert.NoError(t, err)
	route(c)
	assert.Equal(t, http.StatusOK, res.Result().StatusCode)
	content, err := ioutil.ReadAll(res.Result().Body)
	assert.NoError(t, err)

	if len(content) > 0 {
		err = json.Unmarshal(content, response)
		assert.NoError(t, err)
	}
}

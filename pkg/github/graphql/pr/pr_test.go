package pr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"

	"mash/db"
	"mash/pkg/auth"
	service_auth "mash/pkg/auth/service"
	"mash/pkg/change"
	db_change "mash/pkg/change/db"
	graphql_change "mash/pkg/change/graphql"
	service_change "mash/pkg/change/service"
	workers_ci "mash/pkg/ci/workers"
	"mash/pkg/codebase"
	db_codebase "mash/pkg/codebase/db"
	service_codebase "mash/pkg/codebase/service"
	db_comments "mash/pkg/comments/db"
	graphql_comments "mash/pkg/comments/graphql"
	service_comments "mash/pkg/comments/service"
	"mash/pkg/github"
	"mash/pkg/github/client"
	"mash/pkg/github/config"
	db_github "mash/pkg/github/db"
	"mash/pkg/github/routes"
	service_github "mash/pkg/github/service"
	workers_github "mash/pkg/github/workers"
	gqlerrors "mash/pkg/graphql/errors"
	"mash/pkg/graphql/resolvers"
	"mash/pkg/internal/inmemory"
	"mash/pkg/internal/sturdytest"
	"mash/pkg/notification/sender"
	"mash/pkg/posthog"
	"mash/pkg/queue"
	db_review "mash/pkg/review/db"
	db_snapshots "mash/pkg/snapshots/db"
	"mash/pkg/snapshots/snapshotter"
	worker_snapshots "mash/pkg/snapshots/worker"
	db_statuses "mash/pkg/statuses/db"
	graphql_statuses "mash/pkg/statuses/graphql"
	service_statuses "mash/pkg/statuses/service"
	db_suggestion "mash/pkg/suggestions/db"
	service_suggestion "mash/pkg/suggestions/service"
	service_sync "mash/pkg/sync/service"
	"mash/pkg/unidiff"
	service_user "mash/pkg/user/service"
	"mash/pkg/view"
	"mash/pkg/view/events"
	graphql_view "mash/pkg/view/graphql"
	db_activity "mash/pkg/workspace/activity/db"
	activity_sender "mash/pkg/workspace/activity/sender"
	service_activity "mash/pkg/workspace/activity/service"
	graphql_workspace "mash/pkg/workspace/graphql"
	service_workspace "mash/pkg/workspace/service"
	db_workspace_watchers "mash/pkg/workspace/watchers/db"
	graphql_workspace_watchers "mash/pkg/workspace/watchers/graphql"
	service_workspace_watchers "mash/pkg/workspace/watchers/service"
	"mash/vcs"
	"mash/vcs/executor"
	"mash/vcs/testutil"

	"github.com/gin-gonic/gin"
	gh "github.com/google/go-github/v39/github"
	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var allFilesAllowed, _ = unidiff.NewAllower("*")

func TestPRHighLevel(t *testing.T) {
	if os.Getenv("E2E_TEST") == "" {
		t.SkipNow()
	}

	repoProvider := testutil.TestingRepoProvider(t)

	d, err := db.Setup(
		sturdytest.PsqlDbSourceForTesting(),
		true,
		"file://../../../../db/migrations",
	)
	if err != nil {
		panic(err)
	}

	logger, _ := zap.NewDevelopment()
	viewRepo := inmemory.NewInMemoryViewRepo()
	viewUpdates := events.NewInMemory()
	userRepo := inmemory.NewInMemoryUserRepo()
	codebaseRepo := db_codebase.NewRepo(d)
	codebaseUserRepo := inmemory.NewInMemoryCodebaseUserRepo()
	workspaceRepo := inmemory.NewInMemoryWorkspaceRepo()
	gitHubUserRepo := inmemory.NewInMemoryGitHubUserRepo()
	gitHubPRRepo := db_github.NewGitHubPRRepo(d)
	gitHubInstallationRepo := inmemory.NewInMemoryGitHubInstallationRepository()
	gitHubRepositoryRepo := inmemory.NewInMemoryGitHubRepositoryRepo()
	changeRepo := db_change.NewRepo(d)
	changeCommitRepo := db_change.NewCommitRepository(d)
	commentsRepo := db_comments.NewRepo(d)
	snapshotsRepo := db_snapshots.NewRepo(d)
	executorProvider := executor.NewProvider(logger, repoProvider)
	workspaceActivityRepo := db_activity.NewActivityRepo(d)
	workspaceActivityReadsRepo := db_activity.NewActivityReadsRepo(d)
	reviewsRepo := db_review.NewReviewRepository(d)
	eventsSender := events.NewSender(codebaseUserRepo, workspaceRepo, viewUpdates)
	gitSnapshotter := snapshotter.NewGitSnapshotter(snapshotsRepo, workspaceRepo, workspaceRepo, viewRepo, eventsSender, executorProvider, logger)
	snapshotPublisher := worker_snapshots.NewSync(gitSnapshotter)
	activityService := service_activity.New(workspaceActivityReadsRepo, eventsSender)
	activitySender := activity_sender.NewActivitySender(codebaseUserRepo, workspaceActivityRepo, activityService, eventsSender)
	suggestionRepo := db_suggestion.New(d)
	notificationSender := sender.NewNoopNotificationSender()
	postHogClient := posthog.NewFakeClient()
	commentsService := service_comments.New(commentsRepo)

	queue := queue.NewNoop()
	buildQueue := workers_ci.New(zap.NewNop(), queue, nil)
	codebaseService := service_codebase.New(codebaseRepo, codebaseUserRepo)
	userService := service_user.New(zap.NewNop(), userRepo, nil /*jwtService*/, nil /*onetime*/, nil /*emailsender*/, postHogClient)

	changeService := service_change.New(executorProvider, nil, nil, userRepo, changeRepo, changeCommitRepo, nil)
	importer := service_github.ImporterQueue(workers_github.NopImporter())
	cloner := service_github.ClonerQueue(workers_github.NopCloner())
	gitHubService := service_github.New(
		logger,
		gitHubRepositoryRepo,
		gitHubInstallationRepo,
		gitHubUserRepo,
		gitHubPRRepo,
		config.GitHubAppConfig{},
		clientProvider,
		personalClientProvider,
		&importer,
		&cloner,
		workspaceRepo,
		workspaceRepo,
		codebaseUserRepo,
		codebaseRepo,
		executorProvider,
		gitSnapshotter,
		postHogClient,
		nil, // notificationSender
		nil, // eventsSender
		userService,
	)
	workspaceService := service_workspace.New(
		logger,
		postHogClient,

		workspaceRepo,
		workspaceRepo,

		userRepo,
		reviewsRepo,

		commentsService,
		changeService,
		gitHubService,

		activitySender,
		executorProvider,
		eventsSender,
		snapshotPublisher,
		gitSnapshotter,
		buildQueue,
	)

	suggestionsService := service_suggestion.New(
		logger,
		suggestionRepo,
		workspaceService,
		executorProvider,
		gitSnapshotter,
		postHogClient,
		notificationSender,
		eventsSender,
	)

	authService := service_auth.New(codebaseService, userService, workspaceService, nil /*aclProvider*/, nil /*organizationService*/)

	statusesRepo := db_statuses.New(d)
	statusesServcie := service_statuses.New(logger, statusesRepo, eventsSender)
	statusesRootResolver := new(resolvers.StatusesRootResolver)

	syncService := service_sync.New(logger, executorProvider, viewRepo, workspaceRepo, workspaceRepo, gitSnapshotter)

	webhookRoute := routes.Webhook(
		logger,
		config.GitHubAppConfig{},
		posthog.NewFakeClient(),
		gitHubInstallationRepo,
		gitHubRepositoryRepo,
		codebaseRepo,
		executorProvider,
		clientProvider,
		gitHubUserRepo,
		codebaseUserRepo,
		nil,
		gitHubPRRepo,
		workspaceRepo,
		workspaceRepo,
		workspaceService,
		syncService,
		changeRepo,
		changeCommitRepo,
		reviewsRepo,
		eventsSender,
		activitySender,
		statusesServcie,
		commentsService,
		gitHubService,
		buildQueue,
	)

	prResolver := NewResolver(
		logger,
		nil,
		nil,
		statusesRootResolver,
		userRepo,
		codebaseRepo,
		workspaceRepo,
		viewRepo,
		config.GitHubAppConfig{},
		gitHubUserRepo,
		gitHubPRRepo,
		gitHubInstallationRepo,
		gitHubRepositoryRepo,
		clientProvider,
		personalClientProvider,
		viewUpdates,
		posthog.NewFakeClient(),
		authService,
		gitHubService,
	)

	workspaceWatchersRootResolver := new(resolvers.WorkspaceWatcherRootResolver)
	workspaceWatcherRepo := db_workspace_watchers.NewInMemory()
	workspaceWatchersService := service_workspace_watchers.New(workspaceWatcherRepo, eventsSender)

	commentsResolver := graphql_comments.NewResolver(
		userRepo,
		commentsRepo,
		nil,
		workspaceRepo,
		viewRepo,
		codebaseUserRepo,
		changeRepo,
		workspaceWatchersService,
		authService,
		eventsSender,
		viewUpdates,
		sender.NewNoopNotificationSender(),
		activitySender,
		nil,
		nil,
		nil,
		logger,
		posthog.NewFakeClient(),
		executorProvider,
	)

	workspaceResolver := graphql_workspace.NewResolver(
		workspaceRepo,
		codebaseRepo,
		viewRepo,
		nil, // commentRepo
		nil, // snapshotRepo
		nil, // codebaseResolver
		nil, // authorResolver
		nil, // viewResolver
		&commentsResolver,
		nil, // prResolver
		nil, // changeResolver
		nil, // workspaceActivityResolver
		nil, // reviewRootResolver
		nil, // presenseRootResolver
		nil, // suggestitonsRootResolver
		statusesRootResolver,
		workspaceWatchersRootResolver,
		suggestionsService,
		workspaceService,
		authService,
		logger,
		viewUpdates,
		workspaceRepo,
		executorProvider,
		eventsSender,
		gitSnapshotter,
	)

	*workspaceWatchersRootResolver = graphql_workspace_watchers.NewRootResolver(
		logger,
		workspaceWatchersService,
		workspaceService,
		authService,
		viewUpdates,
		nil,
		&workspaceResolver,
	)

	viewResolver := graphql_view.NewResolver(
		viewRepo,
		workspaceRepo,
		nil,
		nil,
		nil,
		nil,
		nil,
		workspaceRepo,
		viewUpdates,
		eventsSender,
		executorProvider,
		logger,
		nil,
		workspaceWatchersService,
		postHogClient,
		nil,
		authService,
	)

	changeResolver := graphql_change.NewResolver(
		changeService,
		changeRepo,
		changeCommitRepo,
		commentsRepo,
		authService,
		&commentsResolver,
		nil, // authorresolver
		statusesRootResolver,
		executorProvider,
		logger,
	)

	*statusesRootResolver = graphql_statuses.New(
		logger,
		statusesServcie,
		changeService,
		workspaceService,
		authService,
		gitHubPRRepo,
		&changeResolver,
		&prResolver,
		viewUpdates,
	)

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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			userID := uuid.NewString()
			viewID := uuid.NewString()
			codebaseID := uuid.NewString()
			codebaseUserID := uuid.NewString()
			gitHubRepositoryID := uuid.NewString()
			gitHubInstallationID := rand.Int63n(50_000_000)
			ctx := auth.NewContext(context.Background(), &auth.Subject{Type: auth.SubjectUser, ID: userID})

			gitHubRepoOwner := uuid.NewString()
			gitHubRepoName := uuid.NewString()

			vw := &view.View{
				ID:         viewID,
				CodebaseID: codebaseID,
				UserID:     userID,
			}
			cu := &codebase.CodebaseUser{
				ID:         codebaseUserID,
				UserID:     userID,
				CodebaseID: codebaseID,
			}
			expT := time.Now().Add(20 * time.Minute)
			ghr := &github.GitHubRepository{
				ID:                               gitHubRepositoryID,
				InstallationID:                   gitHubInstallationID,
				Name:                             gitHubRepoName,
				TrackedBranch:                    "master",
				CodebaseID:                       codebaseID,
				GitHubSourceOfTruth:              true,
				IntegrationEnabled:               true,
				InstallationAccessToken:          str("token"),
				InstallationAccessTokenExpiresAt: &expT,
			}
			ghu := &github.GitHubUser{
				UserID: userID,
			}

			in := &github.GitHubInstallation{
				InstallationID: gitHubInstallationID,
				Owner:          gitHubRepoOwner,
			}

			err = viewRepo.Create(*vw)
			assert.NoError(t, err)
			err = codebaseUserRepo.Create(*cu)
			assert.NoError(t, err)
			err = gitHubRepositoryRepo.Create(*ghr)
			assert.NoError(t, err)
			err = gitHubInstallationRepo.Create(*in)
			assert.NoError(t, err)
			err = gitHubUserRepo.Create(*ghu)
			assert.NoError(t, err)

			err = codebaseRepo.Create(codebase.Codebase{
				ID:              codebaseID,
				ShortCodebaseID: codebase.ShortCodebaseID(codebaseID),
			})
			assert.NoError(t, err)

			// Create GitHub remote
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

			bareRepo, err := vcs.CloneRepoBare(fakeGitHubRemotePath, pathBase)
			if !assert.NoError(t, err) {
				t.FailNow()
			}

			// Create workspace
			wsRes, err := workspaceResolver.CreateWorkspace(ctx, resolvers.CreateWorkspaceArgs{Input: resolvers.CreateWorkspaceInput{CodebaseID: graphql.ID(codebaseID)}})
			assert.NoError(t, err)
			workspaceID := string(wsRes.ID())

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
					DraftDescription: str("<p><em>draft description</em></p>"),
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
					Path:        str("a.txt"),
					LineStart:   i32(1),
					LineEnd:     i32(1),
					LineIsNew:   b(true),
				}})
				assert.NoError(t, err)
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
						PatchIDs:    hunkIDs,
					}},
			)
			if !assert.NoError(t, err) {
				t.Logf("err=%+v", err.(*gqlerrors.SturdyGraphqlError).OriginalError())
			}
			if assert.NotNil(t, createdPullRequestResolver) {
				assert.True(t, createdPullRequestResolver.Open())
				assert.False(t, createdPullRequestResolver.Merged())
			} else {
				t.FailNow()
			}

			// get githubs pull request ID (not the same as the pull request number)
			ghpr, err := gitHubPRRepo.Get(string(createdPullRequestResolver.ID()))
			assert.NoError(t, err)
			gitHubPullRequestID := ghpr.GitHubID

			// PR was closed
			prWebhookEvent(t, "closed", false, userID, gitHubPullRequestID, webhookRoute)

			// Updated PR is closed
			gqlID := graphql.ID(workspaceID)
			updatedPR, err := prResolver.InternalGitHubPullRequestByWorkspaceID(ctx, resolvers.GitHubPullRequestArgs{WorkspaceID: &gqlID})
			assert.NoError(t, err)
			assert.False(t, updatedPR.Open())

			// PR was reopened
			prWebhookEvent(t, "open", false, userID, gitHubPullRequestID, webhookRoute)

			// Updated PR is opened
			gqlID = graphql.ID(workspaceID)
			updatedPR, err = prResolver.InternalGitHubPullRequestByWorkspaceID(ctx, resolvers.GitHubPullRequestArgs{WorkspaceID: &gqlID})
			assert.NoError(t, err)
			assert.True(t, updatedPR.Open())

			// Merge PR
			prWebhookEvent(t, "closed", true, userID, gitHubPullRequestID, webhookRoute)

			// Updated PR is merged
			updatedPR, err = prResolver.InternalGitHubPullRequestByWorkspaceID(ctx, resolvers.GitHubPullRequestArgs{WorkspaceID: &gqlID})
			assert.NoError(t, err)
			assert.False(t, updatedPR.Open())
			assert.True(t, updatedPR.Merged())

			// Rebase or rebase the commit
			if tc.gitHubRebase {
				masterCommit, err := fakeGitHubBareRepo.BranchCommitID("master")
				assert.NoError(t, err)
				branchCommit, err := fakeGitHubBareRepo.BranchCommitID("sturdy-pr-" + workspaceID)
				assert.NoError(t, err)
				newCommit, _, _, err := fakeGitHubBareRepo.CherryPickOnto(branchCommit, masterCommit)
				assert.NoError(t, err)
				err = fakeGitHubBareRepo.MoveBranchToCommit("master", newCommit)
				assert.NoError(t, err)
			} else {
				err = fakeGitHubBareRepo.MergeBranchInto("sturdy-pr-"+workspaceID, "master")
				assert.NoError(t, err)
			}

			// Post-merge push webhook event
			webhookRepoPush := gh.PushEvent{
				Ref:          str("refs/heads/master"),
				Installation: &gh.Installation{ID: &gitHubInstallationID},
			}
			requestWithParams(t, userID, webhookRoute, webhookRepoPush, nil, "push", []gin.Param{})

			// Workspace up to date state is reset after the push event
			ws, err := workspaceRepo.Get(workspaceID)
			assert.NoError(t, err)
			assert.Nil(t, ws.UpToDateWithTrunk)
			assert.Empty(t, ws.DraftDescription) // draft message is reset after push event

			// Imported commit has full metadata
			trunkCommits, err := bareRepo.LogBranch("sturdytrunk", 10)
			assert.NoError(t, err)
			for _, c := range trunkCommits {
				t.Logf("trunkcommit: %+v", c)
			}
			assert.Len(t, trunkCommits, 2)

			// One of the commits is the imported one
			var importIdx = -1
			for idx, c := range trunkCommits {
				if c.Name == "Test Testsson" {
					importIdx = idx
				}
			}

			var changeID change.ID

			if assert.True(t, importIdx >= 0) {
				cc, err := changeCommitRepo.GetByCommitID(trunkCommits[importIdx].ID, codebaseID)
				assert.NoError(t, err)
				changeID = cc.ChangeID
				ch, err := changeRepo.Get(cc.ChangeID)
				assert.NoError(t, err)
				if assert.NotNil(t, ch.Title) {
					assert.Equal(t, "draft description", *ch.Title)
				}
				assert.Equal(t, "<p><em>draft description</em></p>", ch.UpdatedDescription)
			}

			// The workspace should no longer have any comments
			wsResolver, err = workspaceResolver.Workspace(ctx, resolvers.WorkspaceArgs{ID: gqlID})
			assert.NoError(t, err)
			workspaceComments, err = wsResolver.Comments()
			assert.NoError(t, err)
			assert.Len(t, workspaceComments, 0)

			// The new change should have comments
			changeIDgql := graphql.ID(changeID)
			chResolver, err := changeResolver.Change(ctx, resolvers.ChangeArgs{ID: &changeIDgql})
			assert.NoError(t, err)
			if assert.NotNil(t, chResolver) {
				changeComments, err := chResolver.Comments()
				assert.NoError(t, err)
				assert.Len(t, changeComments, 5)
			}
		})
	}
}

func prWebhookEvent(t *testing.T, state string, merged bool, userID string, pullRequestID int64, webhookRoute gin.HandlerFunc) {
	webhookPREvent := gh.PullRequestEvent{
		PullRequest: &gh.PullRequest{
			ID:     &pullRequestID,
			State:  &state,
			Merged: &merged,
		},
	}
	requestWithParams(t, userID, webhookRoute, webhookPREvent, nil, "pull_request", []gin.Param{})
}

func clientProvider(gitHubAppConfig config.GitHubAppConfig, installationID int64) (tokenClient *client.GitHubClients, jwtClient *client.GitHubClients, err error) {
	return &client.GitHubClients{
			Repositories: nil,
			PullRequests: &fakeGitHubPullRequestClient{},
		},
		&client.GitHubClients{
			Apps: &fakeGitHubAppsClient{},
		}, nil
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
		State:  str("open"),
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
		Token:        str("testingtoken"),
		ExpiresAt:    t(time.Now().Add(time.Hour * 3)),
		Permissions:  opts.Permissions,
		Repositories: nil,
	}, nil, nil
}

func requestWithParams(t *testing.T, userID string, route func(*gin.Context), request, response interface{}, reqType string, params []gin.Param) {
	res := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(res)
	c.Params = params

	data, err := json.Marshal(request)
	assert.NoError(t, err)

	c.Request, err = http.NewRequest("GET", "/", bytes.NewReader(data))
	c.Request = c.Request.WithContext(auth.NewContext(context.Background(), &auth.Subject{ID: userID, Type: auth.SubjectUser}))
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

func str(s string) *string {
	return &s
}

func t(t time.Time) *time.Time {
	return &t
}

func b(b bool) *bool {
	return &b
}

func i32(i int32) *int32 {
	return &i
}

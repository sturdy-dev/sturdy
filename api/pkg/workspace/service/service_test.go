package service

import (
	"testing"

	"getsturdy.com/api/pkg/analytics/disabled"
	workers_ci "getsturdy.com/api/pkg/ci/workers"
	"getsturdy.com/api/pkg/events"
	"getsturdy.com/api/pkg/internal/inmemory"
	"getsturdy.com/api/pkg/queue"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"
	"getsturdy.com/api/vcs/provider"
	"getsturdy.com/api/vcs/testutil"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type testCollaborators struct {
	service      Service
	repoProvider provider.RepoProvider
}

func setup(t *testing.T) *testCollaborators {
	logger, _ := zap.NewDevelopment()
	repoProvider := testutil.TestingRepoProvider(t)
	executorProvider := executor.NewProvider(logger, repoProvider)
	workspaceRepo := inmemory.NewInMemoryWorkspaceRepo()
	postHogClient := disabled.NewClient()
	snapshotRepo := inmemory.NewInMemorySnapshotRepo()
	viewRepo := inmemory.NewInMemoryViewRepo()
	viewEvents := events.NewInMemory()
	codebaseUserRepo := inmemory.NewInMemoryCodebaseUserRepo()
	eventsSender := events.NewSender(codebaseUserRepo, workspaceRepo, viewEvents)
	gitSnapshotter := snapshotter.NewGitSnapshotter(snapshotRepo, workspaceRepo, workspaceRepo, viewRepo, eventsSender, executorProvider, logger)
	userRepo := inmemory.NewInMemoryUserRepo()
	queue := queue.NewNoop()
	buildQueue := workers_ci.New(zap.NewNop(), queue, nil)

	service := New(
		logger,
		postHogClient,

		workspaceRepo,
		workspaceRepo,

		userRepo,
		nil, // reviewRepo

		nil, // commentService
		nil, // changeService

		nil, // activitySender
		executorProvider,
		eventsSender,
		nil, // snapshotterQueue
		gitSnapshotter,
		buildQueue,
	)

	return &testCollaborators{
		service,
		repoProvider,
	}
}

func (c *testCollaborators) createCodebase(t *testing.T, id string) vcs.RepoGitWriter {
	repoPath := c.repoProvider.TrunkPath(id)
	repo, err := vcs.CreateBareRepoWithRootCommit(repoPath)
	assert.NoError(t, err)
	return repo
}

func TestCreateNewWorkspace(t *testing.T) {
	c := setup(t)

	request := CreateWorkspaceRequest{
		UserID:     "user-id",
		CodebaseID: "codebase-id",
	}

	c.createCodebase(t, request.CodebaseID)

	ws, err := c.service.Create(request)
	assert.NoError(t, err)

	assert.Equal(t, ws.UserID, request.UserID)
	assert.Equal(t, ws.CodebaseID, request.CodebaseID)
	assert.Equal(t, *ws.Name, "Test Testsson's Workspace")
}

func TestCreateNewWorkspaceWithExplicitName(t *testing.T) {
	c := setup(t)

	request := CreateWorkspaceRequest{
		UserID:     "user-id",
		CodebaseID: "codebase-id",
		Name:       "My New Workspace",
	}

	c.createCodebase(t, request.CodebaseID)

	ws, err := c.service.Create(request)
	assert.NoError(t, err)

	assert.Equal(t, ws.UserID, request.UserID)
	assert.Equal(t, ws.CodebaseID, request.CodebaseID)
	assert.Equal(t, *ws.Name, request.Name)
}

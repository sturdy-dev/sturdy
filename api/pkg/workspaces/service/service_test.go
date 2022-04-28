package service_test

import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"

	"getsturdy.com/api/pkg/codebases"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	service_codebase "getsturdy.com/api/pkg/codebases/service"
	"getsturdy.com/api/pkg/configuration"
	"getsturdy.com/api/pkg/di"
	db_installations "getsturdy.com/api/pkg/installations/db"
	"getsturdy.com/api/pkg/logger"
	module_queue "getsturdy.com/api/pkg/queue/module"
	"getsturdy.com/api/pkg/snapshots"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	service_snapshots "getsturdy.com/api/pkg/snapshots/service"
	db_statuses "getsturdy.com/api/pkg/statuses/db"
	db_suggestions "getsturdy.com/api/pkg/suggestions/db"
	"getsturdy.com/api/pkg/users"
	db_view "getsturdy.com/api/pkg/view/db"
	service_view "getsturdy.com/api/pkg/view/service"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"
	"getsturdy.com/api/vcs/testutil"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func testModule(t *testing.T) di.Module {
	return func(c *di.Container) {
		c.Import(service_snapshots.Module)
		c.Import(service_codebase.Module)
		c.Import(service_workspace.Module)
		c.Import(service_view.Module)

		c.ImportWithForce(db_snapshots.TestModule)
		c.ImportWithForce(db_view.TestModule)
		c.ImportWithForce(db_workspaces.TestModule)
		c.ImportWithForce(db_suggestions.TestModule)
		c.ImportWithForce(db_codebases.TestModule)
		c.ImportWithForce(db_installations.TestModule)
		c.ImportWithForce(db_statuses.TestModule)
		c.ImportWithForce(module_queue.TestModule)
		c.ImportWithForce(configuration.TestModule)
		c.RegisterWithForce(logger.NewTest)

		c.RegisterWithForce(func() *sqlx.DB { return nil }) // make sure db is not used
		c.Register(func() *testing.T { return t })
		c.RegisterWithForce(testutil.TestingRepoProvider)
	}
}

type testCase struct {
	snapshotService  *service_snapshots.Service
	workspaceService *service_workspace.Service
	codebaseService  *service_codebase.Service
	viewService      *service_view.Service
	executorProvider executor.Provider

	userID     users.ID
	codebaseID codebases.ID
}

func setup(t *testing.T) *testCase {
	tc := &testCase{}
	if !assert.NoError(t, di.Init(testModule(t)).To(
		&tc.snapshotService, &tc.workspaceService, &tc.codebaseService, &tc.viewService, &tc.executorProvider,
	)) {
		t.FailNow()
	}

	userID := users.ID(uuid.NewString())
	tc.userID = userID

	ctx := context.Background()

	cb, err := tc.codebaseService.Create(ctx, userID, "test", nil)
	assert.NoError(t, err)
	tc.codebaseID = cb.ID

	return tc
}

func TestCreate_no_name(t *testing.T) {
	tc := setup(t)

	ctx := context.Background()

	ws, err := tc.workspaceService.Create(ctx, service_workspace.CreateWorkspaceRequest{
		UserID:     tc.userID,
		CodebaseID: tc.codebaseID,
	})
	assert.NoError(t, err)
	assert.Equal(t, "Untitled draft", *ws.Name)
}

func TestCreate_with_name(t *testing.T) {
	tc := setup(t)

	ctx := context.Background()

	ws, err := tc.workspaceService.Create(ctx, service_workspace.CreateWorkspaceRequest{
		UserID:     tc.userID,
		CodebaseID: tc.codebaseID,
		Name:       "test",
	})
	assert.NoError(t, err)
	assert.Equal(t, "test", *ws.Name)
}

func TestWorkspace_SetSnapshot(t *testing.T) {
	tc := setup(t)

	ctx := context.Background()

	ws, err := tc.workspaceService.Create(ctx, service_workspace.CreateWorkspaceRequest{UserID: tc.userID, CodebaseID: tc.codebaseID})
	assert.NoError(t, err)

	vw, err := tc.viewService.Create(ctx, tc.userID, ws, nil, nil)
	assert.NoError(t, err)

	assert.NoError(t, tc.executorProvider.New().Write(writeFile("test.txt", []byte("test"))).ExecView(tc.codebaseID, vw.ID, "make some changes"))

	firstSnapshot, err := tc.snapshotService.Snapshot(tc.codebaseID, ws.ID, snapshots.Action("testing"), service_snapshots.WithOnView(vw.ID))
	assert.NoError(t, err)

	assert.NoError(t, tc.executorProvider.New().Write(writeFile("test.txt", []byte("test2"))).ExecView(tc.codebaseID, vw.ID, "make more changes"))

	secondSnapshot, err := tc.snapshotService.Snapshot(tc.codebaseID, ws.ID, snapshots.Action("testing"), service_snapshots.WithOnView(vw.ID))
	assert.NoError(t, err)
	assert.Equal(t, firstSnapshot.ID, *secondSnapshot.PreviousSnapshotID)

	wsBeforeUndo, err := tc.workspaceService.GetByID(ctx, ws.ID)
	assert.NoError(t, err)
	assert.Equal(t, secondSnapshot.ID, *wsBeforeUndo.LatestSnapshotID)
	assert.NoError(t, tc.executorProvider.New().Read(verifySnapshot(secondSnapshot)).ExecView(tc.codebaseID, vw.ID, "verify snapshot"))

	assert.NoError(t, tc.workspaceService.SetSnapshot(ctx, ws, firstSnapshot))

	wsAfterUndo, err := tc.workspaceService.GetByID(ctx, ws.ID)
	assert.NoError(t, err)
	assert.Equal(t, firstSnapshot.ID, *wsAfterUndo.LatestSnapshotID)
	assert.NoError(t, tc.executorProvider.New().Read(verifySnapshot(firstSnapshot)).ExecView(tc.codebaseID, vw.ID, "verify snapshot"))

	assert.NoError(t, tc.workspaceService.SetSnapshot(ctx, ws, secondSnapshot))

	wsAfterRedo, err := tc.workspaceService.GetByID(ctx, ws.ID)
	assert.NoError(t, err)
	assert.Equal(t, secondSnapshot.ID, *wsAfterRedo.LatestSnapshotID)
	assert.NoError(t, tc.executorProvider.New().Read(verifySnapshot(secondSnapshot)).ExecView(tc.codebaseID, vw.ID, "verify snapshot"))
}

func verifySnapshot(snapshot *snapshots.Snapshot) func(vcs.RepoReader) error {
	return func(repo vcs.RepoReader) error {
		head, err := repo.HeadCommit()
		if err != nil {
			return fmt.Errorf("could not get head commit: %w", err)
		}
		defer head.Free()

		snapshotParents, err := repo.GetCommitParents(snapshot.CommitSHA)
		if err != nil {
			return fmt.Errorf("could not get snapshot parents: %w", err)
		}

		if len(snapshotParents) != 1 {
			return fmt.Errorf("expected snapshot to have 1 parent, got %d", len(snapshotParents))
		}

		if head.Id().String() != snapshotParents[0] {
			return fmt.Errorf("expected commit %s, got %s", snapshotParents[0], head.Id().String())
		}

		return nil
	}
}

func writeFile(filename string, content []byte) func(vcs.RepoWriter) error {
	return func(repo vcs.RepoWriter) error {
		file, err := os.Create(path.Join(repo.Path(), filename))
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err := file.Write(content); err != nil {
			return err
		}
		return nil
	}
}

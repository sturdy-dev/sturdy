package service_test

import (
	"context"
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
		c.ImportWithForce(module_queue.TestModule)
		c.ImportWithForce(configuration.TestModule)
		c.RegisterWithForce(logger.NewTest)

		c.RegisterWithForce(func() *sqlx.DB { return nil })
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

	userID      users.ID
	codebaseID  codebases.ID
	workspaceID string
	viewID      string
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

	ws, err := tc.workspaceService.Create(ctx, service_workspace.CreateWorkspaceRequest{UserID: userID, CodebaseID: cb.ID})
	assert.NoError(t, err)
	tc.workspaceID = ws.ID

	vw, err := tc.viewService.Create(ctx, userID, ws, nil, nil)
	assert.NoError(t, err)
	tc.viewID = vw.ID

	return tc
}

func TestSnapshot_noChanges(t *testing.T) {
	tc := setup(t)

	snapOnce, err := tc.snapshotService.Snapshot(tc.codebaseID, tc.workspaceID, snapshots.ActionViewSync, service_snapshots.WithOnView(tc.viewID))
	assert.NoError(t, err)

	snapTwice, err := tc.snapshotService.Snapshot(tc.codebaseID, tc.workspaceID, snapshots.ActionViewSync, service_snapshots.WithOnView(tc.viewID))
	assert.NoError(t, err)

	assert.Equal(t, snapOnce.ID, snapTwice.ID, "snapshots must be identical, no changes were made")
}

func TestSnapshot_noChangesDeleted(t *testing.T) {
	tc := setup(t)

	snapOnce, err := tc.snapshotService.Snapshot(tc.codebaseID, tc.workspaceID, snapshots.ActionViewSync, service_snapshots.WithOnView(tc.viewID))
	assert.NoError(t, err)

	assert.NoError(t, tc.snapshotService.Delete(context.Background(), snapOnce))

	snapTwice, err := tc.snapshotService.Snapshot(tc.codebaseID, tc.workspaceID, snapshots.ActionViewSync, service_snapshots.WithOnView(tc.viewID))
	assert.NoError(t, err)

	assert.NotEqual(t, snapOnce.ID, snapTwice.ID, "snapshots must not be identical, the first one was removed")
}

func TestSnapshot_changes(t *testing.T) {
	tc := setup(t)

	snapOnce, err := tc.snapshotService.Snapshot(tc.codebaseID, tc.workspaceID, snapshots.ActionViewSync, service_snapshots.WithOnView(tc.viewID))
	assert.NoError(t, err)

	assert.NoError(t, tc.executorProvider.New().Write(writeFile("test.txt", []byte("test"))).ExecView(tc.codebaseID, tc.viewID, "make some changes"))

	snapTwice, err := tc.snapshotService.Snapshot(tc.codebaseID, tc.workspaceID, snapshots.ActionViewSync, service_snapshots.WithOnView(tc.viewID))
	assert.NoError(t, err)

	assert.NotEqual(t, snapOnce.ID, snapTwice.ID, "snapshots must not be identical, changes were made")
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

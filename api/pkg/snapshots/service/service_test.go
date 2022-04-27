package service_test

import (
	"context"
	"testing"

	"getsturdy.com/api/pkg/codebases"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	service_codebase "getsturdy.com/api/pkg/codebases/service"
	"getsturdy.com/api/pkg/configuration"
	"getsturdy.com/api/pkg/di"
	db_installations "getsturdy.com/api/pkg/installations/db"
	"getsturdy.com/api/pkg/logger"
	module_queue "getsturdy.com/api/pkg/queue/module"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	service_snapshots "getsturdy.com/api/pkg/snapshots/service"
	db_suggestions "getsturdy.com/api/pkg/suggestions/db"
	"getsturdy.com/api/pkg/users"
	db_view "getsturdy.com/api/pkg/view/db"
	service_view "getsturdy.com/api/pkg/view/service"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"
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

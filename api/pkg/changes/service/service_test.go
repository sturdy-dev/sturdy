package service_test

import (
	"context"
	"os/exec"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
	"go.uber.org/zap"

	module_api "getsturdy.com/api/pkg/api/module"
	"getsturdy.com/api/pkg/changes/service"
	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/configuration"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/internal/dbtest"
	queue "getsturdy.com/api/pkg/queue/module"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"
	"getsturdy.com/api/vcs/provider"
)

func module(t *testing.T) di.Module {
	return func(c *di.Container) {
		// TODO: reduce scope
		c.Import(module_api.Module)
		c.ImportWithForce(configuration.TestModule)
		c.ImportWithForce(queue.TestModule)
		c.Register(func() *testing.T { return t })
		c.RegisterWithForce(dbtest.DB)
	}
}

func TestChangelog(t *testing.T) {
	type deps struct {
		dig.In
		ExecutorProvider executor.Provider
		RepoProvider     provider.RepoProvider
		Logger           *zap.Logger
		Service          *service.Service
	}

	var d deps
	if !assert.NoError(t, di.Init(module(t)).To(&d)) {
		t.FailNow()
	}

	svc := d.Service
	codebaseID := codebases.ID(uuid.NewString())

	barePath := d.RepoProvider.TrunkPath(codebaseID)
	_, err := vcs.CreateBareRepoWithRootCommit(barePath)
	assert.NoError(t, err)

	viewID := uuid.NewString()
	viewPath := d.RepoProvider.ViewPath(codebaseID, viewID)

	t.Log("viewPath", viewPath)

	_, err = vcs.CloneRepo(barePath, viewPath)
	assert.NoError(t, err)

	output, err := exec.Command("bash", "testdata/generate-repo.sh", viewPath).CombinedOutput()
	if !assert.NoError(t, err) {
		t.Logf("output: %s", string(output))
	}

	err = d.ExecutorProvider.New().Write(func(writer vcs.RepoWriter) error {
		err = writer.ForcePush(d.Logger, "sturdytrunk")
		assert.NoError(t, err)
		return nil
	}).AllowRebasingState().ExecView(codebaseID, viewID, "test")
	assert.NoError(t, err)

	ctx := context.Background()

	log1, err := svc.Changelog(ctx, codebaseID, 3, nil)
	assert.NoError(t, err)

	expected := []struct {
		message string
	}{
		{"Merge ws2"},
		{"Merge ws1"},
		{"6"},
		{"3"},
		{"2"},
		{"1"},
	}

	if assert.Len(t, log1, 3) {
		for k, v := range log1 {
			assert.Equal(t, expected[k].message, v.UpdatedDescription, "pos=%d", k)
		}
	}

	log2, err := svc.Changelog(ctx, codebaseID, 10, &log1[len(log1)-1].ID)
	assert.NoError(t, err)

	if assert.Len(t, log2, 3) {
		for k, v := range log2 {
			assert.Equal(t, expected[k+3].message, v.UpdatedDescription, "pos=%d", k)
		}
	}
}

func Test_Parent_Child_navigation(t *testing.T) {
	type deps struct {
		dig.In
		ExecutorProvider executor.Provider
		RepoProvider     provider.RepoProvider
		Logger           *zap.Logger
		Service          *service.Service
	}

	var d deps
	if !assert.NoError(t, di.Init(module(t)).To(&d)) {
		t.FailNow()
	}

	svc := d.Service
	codebaseID := codebases.ID(uuid.NewString())

	barePath := d.RepoProvider.TrunkPath(codebaseID)
	_, err := vcs.CreateBareRepoWithRootCommit(barePath)
	assert.NoError(t, err)

	viewID := uuid.NewString()
	viewPath := d.RepoProvider.ViewPath(codebaseID, viewID)

	t.Log("viewPath", viewPath)

	_, err = vcs.CloneRepo(barePath, viewPath)
	assert.NoError(t, err)

	output, err := exec.Command("bash", "testdata/generate-repo.sh", viewPath).CombinedOutput()
	if !assert.NoError(t, err) {
		t.Logf("output: %s", string(output))
	}

	assert.NoError(t, d.ExecutorProvider.New().Write(func(writer vcs.RepoWriter) error {
		assert.NoError(t, writer.ForcePush(d.Logger, "sturdytrunk"))
		return nil
	}).AllowRebasingState().ExecView(codebaseID, viewID, "test"))

	ctx := context.Background()

	changes, err := svc.Changelog(ctx, codebaseID, 100, nil)
	assert.NoError(t, err)

	if !assert.Len(t, changes, 6) {
		t.FailNow()
	}

	secondChange := changes[1]

	secondChangeParent, err := svc.ParentChange(ctx, secondChange)
	assert.NoError(t, err)
	assert.Equal(t, secondChangeParent.ID, *secondChange.ParentChangeID)
	assert.Equal(t, secondChangeParent.ID, changes[2].ID)

	secondChangeChild, err := svc.ChildChange(ctx, secondChange)
	assert.NoError(t, err)
	assert.Equal(t, *secondChangeChild.ParentChangeID, secondChange.ID)
	assert.Equal(t, secondChangeChild.ID, changes[0].ID)

	lastChange := changes[len(changes)-1]
	_, lastChangeParentErr := svc.ParentChange(ctx, lastChange)
	assert.Equal(t, lastChangeParentErr, service.ErrNotFound)

	firstChange := changes[0]
	_, firstChangeChildErr := svc.ChildChange(ctx, firstChange)
	assert.Equal(t, firstChangeChildErr, service.ErrNotFound)
}

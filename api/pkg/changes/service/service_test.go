package service_test

import (
	"context"
	"os/exec"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
	"go.uber.org/zap"

	"getsturdy.com/api/pkg/changes/service"
	module_configuration "getsturdy.com/api/pkg/configuration/module"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/installations"
	"getsturdy.com/api/pkg/internal/inmemory"
	module_logger "getsturdy.com/api/pkg/logger/module"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"
	module_vcs "getsturdy.com/api/vcs/module"
	"getsturdy.com/api/vcs/provider"
)

func module(c *di.Container) {
	ctx := context.Background()
	c.Register(func() context.Context {
		return ctx
	})

	c.Register(func() *installations.Installation {
		return &installations.Installation{ID: uuid.NewString()}
	})

	c.Import(module_logger.Module)
	c.Import(module_vcs.Module)
	c.Import(module_configuration.TestingModule)
}

func TestChangelog(t *testing.T) {
	changeRepo := inmemory.NewInMemoryChangeRepo()

	type deps struct {
		dig.In
		ExecutorProvider executor.Provider
		RepoProvider     provider.RepoProvider
		Logger           *zap.Logger
	}

	var d deps
	if !assert.NoError(t, di.Init(&d, module)) {
		t.FailNow()
	}

	codebaseID := uuid.NewString()

	svc := service.New(changeRepo, d.Logger, d.ExecutorProvider)

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

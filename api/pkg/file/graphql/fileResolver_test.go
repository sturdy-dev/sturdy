package graphql

import (
	"context"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"
	"getsturdy.com/api/vcs/testutil"
	"io/ioutil"
	"path"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestFileResolver(t *testing.T) {
	repoProvider := testutil.TestingRepoProvider(t)
	logger := zap.NewNop()
	executorProvider := executor.NewProvider(logger, repoProvider)

	codebaseID := uuid.NewString()
	viewID := uuid.NewString()
	trunkPath := repoProvider.TrunkPath(codebaseID)
	_, err := vcs.CreateBareRepoWithRootCommit(trunkPath)
	assert.NoError(t, err)
	viewPath := repoProvider.ViewPath(codebaseID, viewID)
	viewRepo, err := vcs.CloneRepo(trunkPath, viewPath)
	assert.NoError(t, err)

	_ = viewRepo

	root := NewFileRootResolver(executorProvider)
	fileResolver, err := root.InternalFile(context.Background(), codebaseID, "README.md", "README.markdown")
	assert.Error(t, err, gqlerrors.ErrNotFound)
	assert.Nil(t, fileResolver)

	t.Log(trunkPath)

	fileNames := []string{"readme.markdown", "README.md"}
	for _, fileName := range fileNames {
		t.Run(fileName, func(t *testing.T) {
			// Create readme.markdown
			err = ioutil.WriteFile(path.Join(viewPath, fileName), []byte("# hey from "+fileName), 0o644)
			assert.NoError(t, err)
			_, err = viewRepo.AddAndCommit(fileName)
			assert.NoError(t, err)
			err = viewRepo.Push(logger, "sturdytrunk")
			assert.NoError(t, err)

			fileResolver, err = root.InternalFile(context.Background(), codebaseID, "README.md", "README.markdown")
			assert.NoError(t, err)
			if assert.NotNil(t, fileResolver) {
				fileResolver, ok := fileResolver.ToFile()
				assert.True(t, ok)
				assert.Equal(t, fileName, fileResolver.Path())
				assert.Equal(t, "# hey from "+fileName, fileResolver.Contents())
			}
		})
	}
}

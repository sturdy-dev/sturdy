package graphql

import (
	"context"
	"io/ioutil"
	"path"
	"strings"
	"testing"
	"time"

	"getsturdy.com/api/pkg/auth"
	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/codebase/acl"
	provider_acl "getsturdy.com/api/pkg/codebase/acl/provider"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/internal/inmemory"
	service_user "getsturdy.com/api/pkg/users/service"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"
	"getsturdy.com/api/vcs/testutil"

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
	viewGitRepo, err := vcs.CloneRepo(trunkPath, viewPath)
	assert.NoError(t, err)

	aclRepo := inmemory.NewInMemoryAclRepo()
	userRepo := inmemory.NewInMemoryUserRepo()

	userService := service_user.New(
		zap.NewNop(),
		userRepo,
		nil,
		nil,
		nil,
		nil,
	)

	aclProvider := provider_acl.New(
		aclRepo,
		nil,
		nil,
	)

	authService := service_auth.New(
		nil,
		userService,
		nil,
		aclProvider,
		nil,
	)

	aclID := uuid.NewString()
	userID := uuid.NewString()
	singleFileOnlyUserID := uuid.NewString()

	rawPolicy := strings.NewReplacer(
		"__USER_ID__", userID,
		"__SINGLE_FILE_ONLY_USER_ID__", singleFileOnlyUserID,
		"__CODEBASE_ID__", codebaseID,
	).Replace(`{
  "groups": [
    {
      "id": "our-user",
      "members": ["users::__USER_ID__"],
    },
  ],
  "tests": [
    {
      "id": "can manage access control",
      "principal": "users::__USER_ID__",
      "allow": "write",
      "resource": "acls::__CODEBASE_ID__",
    },
  ],
  "rules": [
    {
      "id": "everyone can manage access control",
      "principals": ["groups::our-user"],
      "action": "write",
      "resources": ["acls::__CODEBASE_ID__"],
    },
    {
      "id": "user can access all but one file",
      "principals": ["users::__USER_ID__"],
      "action": "write",
      "resources": [ "files::*", "files::!not_allowed.txt" ],
    },
    {
      "id": "user can access only readme",
      "principals": ["users::__SINGLE_FILE_ONLY_USER_ID__"],
      "action": "write",
      "resources": [ "files::README.md" ],
    }
  ],
}`)

	rule := acl.ACL{
		ID:         acl.ID(aclID),
		CodebaseID: codebaseID,
		CreatedAt:  time.Now(),
		RawPolicy:  rawPolicy,
	}

	err = aclRepo.Create(context.Background(), rule)
	assert.NoError(t, err)

	ctx := auth.NewContext(context.Background(), &auth.Subject{ID: userID, Type: auth.SubjectUser})

	root := NewFileRootResolver(executorProvider, authService)
	fileResolver, err := root.InternalFile(ctx, codebaseID, "README.md", "README.markdown")
	assert.Error(t, err, gqlerrors.ErrNotFound)
	assert.Nil(t, fileResolver)

	t.Log(trunkPath)

	fileNames := []string{"readme.markdown", "README.md", "not_allowed.txt"}
	for _, fileName := range fileNames {
		// Create and push files
		err = ioutil.WriteFile(path.Join(viewPath, fileName), []byte("# hey from "+fileName), 0o644)
		assert.NoError(t, err)
		_, err = viewGitRepo.AddAndCommit(fileName)
		assert.NoError(t, err)
		err = viewGitRepo.Push(logger, "sturdytrunk")
		assert.NoError(t, err)
	}

	{
		fileResolver, err = root.InternalFile(ctx, codebaseID, "README.md", "README.markdown")
		assert.NoError(t, err)
		if assert.NotNil(t, fileResolver) {
			fileResolver, ok := fileResolver.ToFile()
			assert.True(t, ok)
			assert.Equal(t, "README.md", fileResolver.Path())
			assert.Equal(t, "# hey from README.md", fileResolver.Contents())
		}
	}

	{
		fileResolver, err = root.InternalFile(ctx, codebaseID, "/")
		assert.NoError(t, err)
		if assert.NotNil(t, fileResolver) {
			dir, ok := fileResolver.ToDirectory()
			if assert.True(t, ok) {
				t.Log(dir.Path())

				var names []string
				children, err := dir.Children(ctx)
				assert.NoError(t, err)

				for _, ch := range children {
					names = append(names, ch.Path())
				}

				// note how not_allowed.txt is missing
				assert.Equal(t, []string{"readme.markdown", "README.md"}, names)
			}
		}
	}

	{
		fileResolver, err = root.InternalFile(ctx, codebaseID, "not_allowed.txt")
		assert.Error(t, err, gqlerrors.ErrNotFound)
	}

	{
		singleFileUserCtx := auth.NewContext(context.Background(), &auth.Subject{ID: singleFileOnlyUserID, Type: auth.SubjectUser})

		fileResolver, err = root.InternalFile(singleFileUserCtx, codebaseID, "/")
		assert.NoError(t, err)
		if assert.NotNil(t, fileResolver) {
			dir, ok := fileResolver.ToDirectory()
			if assert.True(t, ok) {
				t.Log(dir.Path())

				var names []string
				children, err := dir.Children(singleFileUserCtx)
				assert.NoError(t, err)

				for _, ch := range children {
					names = append(names, ch.Path())
				}

				// only get README.md
				assert.Equal(t, []string{"README.md"}, names)
			}
		}

	}
}

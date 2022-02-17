package graphql

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"getsturdy.com/api/pkg/auth"
	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/codebase"
	"getsturdy.com/api/pkg/codebase/acl"
	provider_acl "getsturdy.com/api/pkg/codebase/acl/provider"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/internal/inmemory"
	db_users "getsturdy.com/api/pkg/users/db"
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
	userRepo := db_users.NewMemory()

	userService := service_user.New(
		zap.NewNop(),
		userRepo,
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
	restrictedUserReadmeUserID := uuid.NewString()
	restrictedUserDirUserID := uuid.NewString()
	restrictedUserDirNoTrailingSlashUserID := uuid.NewString()
	restrictedUserDirFooUserID := uuid.NewString()

	rawPolicy := strings.NewReplacer(
		"__USER_ID__", userID,
		"__RESTRICTED_USER_README_USER_ID__", restrictedUserReadmeUserID,
		"__RESTRICTED_USER_DIR_USER_ID__", restrictedUserDirUserID,
		"__RESTRICTED_USER_DIR_NO_TRAILING_SLASH_USER_ID__", restrictedUserDirNoTrailingSlashUserID,
		"__RESTRICTED_USER_DIR_FOO_USER_ID__", restrictedUserDirFooUserID,
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
      "id": "user __USER_ID__",
      "principals": ["users::__USER_ID__"],
      "action": "write",
      "resources": [ "files::*", "files::!not_allowed.txt" ],
    },
    {
      "id": "user __RESTRICTED_USER_README_USER_ID__",
      "principals": ["users::__RESTRICTED_USER_README_USER_ID__"],
      "action": "write",
      "resources": [ "files::README.md"],
    },
    {
      "id": "user __RESTRICTED_USER_DIR_USER_ID__",
      "principals": ["users::__RESTRICTED_USER_DIR_USER_ID__"],
      "action": "write",
      "resources": [ "files::/pkg/**/*", "files::/pkg/", "files::/dir/**/*", "files::/dir/" ],
    },
    {
      "id": "user __RESTRICTED_USER_DIR_NO_TRAILING_SLASH_USER_ID__",
      "principals": ["users::__RESTRICTED_USER_DIR_NO_TRAILING_SLASH_USER_ID__"],
      "action": "write",
      "resources": [ "files::/pkg/**/*", "files::/pkg/", "files::/dir/**/*", "files::/dir" ],
    },
    {
      "id": "user __RESTRICTED_USER_DIR_USER_ID__",
      "principals": ["users::__RESTRICTED_USER_DIR_USER_ID__"],
      "action": "write",
      "resources": [ "files::/pkg/**/*", "files::/pkg/", "files::/dir/foo/**/*", "files::/dir/foo/" ],
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
	fileResolver, err := root.InternalFile(ctx, &codebase.Codebase{ID: codebaseID}, "README.md", "README.markdown")
	assert.Error(t, err, gqlerrors.ErrNotFound)
	assert.Nil(t, fileResolver)

	fileNames := []string{"readme.markdown", "README.md", "not_allowed.txt", "dir/foo/bar/a.txt", "dir/foo/bar/b.txt"}
	for _, fileName := range fileNames {
		// Create and push files
		err = os.MkdirAll(path.Dir(path.Join(viewPath, fileName)), 0o777)
		assert.NoError(t, err)
		err = ioutil.WriteFile(path.Join(viewPath, fileName), []byte("# hey from "+fileName), 0o644)
		assert.NoError(t, err)
		_, err = viewGitRepo.AddAndCommit(fileName)
		assert.NoError(t, err)
		err = viewGitRepo.Push(logger, "sturdytrunk")
		assert.NoError(t, err)
	}

	t.Run("get readme", func(t *testing.T) {
		fileResolver, err = root.InternalFile(ctx, &codebase.Codebase{ID: codebaseID}, "README.md", "README.markdown")
		assert.NoError(t, err)
		if assert.NotNil(t, fileResolver) {
			fileResolver, ok := fileResolver.ToFile()
			assert.True(t, ok)
			assert.Equal(t, "README.md", fileResolver.Path())
			assert.Equal(t, "# hey from README.md", fileResolver.Contents())
		}

	})

	t.Run("list root not see not_allowed.txt", func(t *testing.T) {
		fileResolver, err = root.InternalFile(ctx, &codebase.Codebase{ID: codebaseID}, "/")
		assert.NoError(t, err)
		if assert.NotNil(t, fileResolver) {
			dir, ok := fileResolver.ToDirectory()
			if assert.True(t, ok) {
				var names []string
				children, err := dir.Children(ctx)
				assert.NoError(t, err)

				for _, ch := range children {
					names = append(names, ch.Path())
				}

				// note how not_allowed.txt is missing
				assert.Equal(t, []string{"dir", "readme.markdown", "README.md"}, names)
			}
		}
	})

	t.Run("list not_allowed.txt not allowed", func(t *testing.T) {
		fileResolver, err = root.InternalFile(ctx, &codebase.Codebase{ID: codebaseID}, "not_allowed.txt")
		assert.Error(t, err, gqlerrors.ErrNotFound)
	})

	t.Run("as restrictedUserDirUserID", func(t *testing.T) {
		restrictedUserCtx := auth.NewContext(context.Background(), &auth.Subject{ID: restrictedUserDirUserID, Type: auth.SubjectUser})

		{
			fileResolver, err = root.InternalFile(restrictedUserCtx, &codebase.Codebase{ID: codebaseID}, "/")
			assert.NoError(t, err)
			if assert.NotNil(t, fileResolver) {
				dir, ok := fileResolver.ToDirectory()
				if assert.True(t, ok) {
					var names []string
					children, err := dir.Children(restrictedUserCtx)
					assert.NoError(t, err)

					for _, ch := range children {
						names = append(names, ch.Path())
					}

					// TODO: To match the mutagen implementation, this user should see dir
					// assert.Equal(t, []string{"dir"}, names)
					assert.Nil(t, names)
				}
			}
		}

		{
			fileResolver, err = root.InternalFile(restrictedUserCtx, &codebase.Codebase{ID: codebaseID}, "/dir")
			assert.NoError(t, err)
			if assert.NotNil(t, fileResolver) {
				dir, ok := fileResolver.ToDirectory()
				if assert.True(t, ok) {
					var names []string
					children, err := dir.Children(restrictedUserCtx)
					assert.NoError(t, err)
					for _, ch := range children {
						names = append(names, ch.Path())
					}
					assert.Equal(t, []string{"dir/foo"}, names)
				}
			}
		}
	})

	t.Run("as restrictedUserReadmeUserID", func(t *testing.T) {
		restrictedUserCtx := auth.NewContext(context.Background(), &auth.Subject{ID: restrictedUserReadmeUserID, Type: auth.SubjectUser})
		{
			fileResolver, err = root.InternalFile(restrictedUserCtx, &codebase.Codebase{ID: codebaseID}, "/")
			assert.NoError(t, err)
			if assert.NotNil(t, fileResolver) {
				dir, ok := fileResolver.ToDirectory()
				if assert.True(t, ok) {
					var names []string
					children, err := dir.Children(restrictedUserCtx)
					assert.NoError(t, err)
					for _, ch := range children {
						names = append(names, ch.Path())
					}
					assert.Equal(t, []string{"README.md"}, names)
				}
			}
		}
	})

	t.Run("as restrictedUserDirNoTrailingSlashCtx", func(t *testing.T) {
		restrictedUserCtx := auth.NewContext(context.Background(), &auth.Subject{ID: restrictedUserDirNoTrailingSlashUserID, Type: auth.SubjectUser})
		{
			fileResolver, err = root.InternalFile(restrictedUserCtx, &codebase.Codebase{ID: codebaseID}, "/")
			assert.NoError(t, err)
			if assert.NotNil(t, fileResolver) {
				dir, ok := fileResolver.ToDirectory()
				if assert.True(t, ok) {
					var names []string
					children, err := dir.Children(restrictedUserCtx)
					assert.NoError(t, err)

					for _, ch := range children {
						names = append(names, ch.Path())
					}

					// only get dir
					assert.Equal(t, []string{"dir"}, names)
				}
			}
		}
	})

	t.Run("as restrictedUserDirFooCtx", func(t *testing.T) {
		restrictedUserCtx := auth.NewContext(context.Background(), &auth.Subject{ID: restrictedUserDirFooUserID, Type: auth.SubjectUser})
		{
			fileResolver, err = root.InternalFile(restrictedUserCtx, &codebase.Codebase{ID: codebaseID}, "/")
			assert.NoError(t, err)
			if assert.NotNil(t, fileResolver) {
				dir, ok := fileResolver.ToDirectory()
				if assert.True(t, ok) {
					var names []string
					children, err := dir.Children(restrictedUserCtx)
					assert.NoError(t, err)

					for _, ch := range children {
						names = append(names, ch.Path())
					}

					// TODO: To match the mutagen implementation, this user should see dir
					// assert.Equal(t, []string{"dir"}, names)
					assert.Nil(t, names)
				}
			}
		}
	})
}

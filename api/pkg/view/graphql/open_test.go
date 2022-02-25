package graphql_test

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"go.uber.org/dig"

	module_api "getsturdy.com/api/pkg/api/module"
	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/codebase"
	db_codebase "getsturdy.com/api/pkg/codebase/db"
	module_configuration "getsturdy.com/api/pkg/configuration/module"
	"getsturdy.com/api/pkg/di"
	module_github "getsturdy.com/api/pkg/github/module"
	"getsturdy.com/api/pkg/graphql/resolvers"
	module_snapshots "getsturdy.com/api/pkg/snapshots/module"
	"getsturdy.com/api/pkg/users"
	db_user "getsturdy.com/api/pkg/users/db"
	"getsturdy.com/api/pkg/view"
	db_view "getsturdy.com/api/pkg/view/db"
	"getsturdy.com/api/pkg/workspaces"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/provider"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func module(c *di.Container) {
	ctx := context.Background()
	c.Register(func() context.Context {
		return ctx
	})

	c.Import(module_api.Module)
	c.Import(module_configuration.TestingModule)
	c.Import(module_snapshots.TestingModule)

	// OSS version
	c.Import(module_github.Module)
}

func TestUpdateViewWorkspace(t *testing.T) {
	if os.Getenv("E2E_TEST") == "" {
		t.SkipNow()
	}

	type deps struct {
		dig.In
		ViewRootResolver resolvers.ViewRootResolver

		RepoProvider provider.RepoProvider

		UserRepo         db_user.Repository
		CodebaseRepo     db_codebase.CodebaseRepository
		WorkspaceRepo    db_workspaces.Repository
		ViewRepo         db_view.Repository
		CodebaseUserRepo db_codebase.CodebaseUserRepository
	}

	var d deps
	if !assert.NoError(t, di.Init(&d, module)) {
		t.FailNow()
	}

	userID := users.ID(uuid.NewString())
	err := d.UserRepo.Create(&users.User{ID: userID, Email: userID.String() + "@test.com"})
	assert.NoError(t, err)

	viewResolver := d.ViewRootResolver
	repoProvider := d.RepoProvider
	codebaseRepo := d.CodebaseRepo
	workspaceRepo := d.WorkspaceRepo
	codebaseUserRepo := d.CodebaseUserRepo
	viewRepo := d.ViewRepo

	authCtx := auth.NewContext(context.Background(), &auth.Subject{Type: auth.SubjectUser, ID: userID.String()})

	type steps struct {
		workspace       string
		expected        string
		toWrite         string
		toWriteNewFile  string
		expectInNewFile string
	}

	cases := []struct {
		name  string
		steps []steps
	}{
		{
			name: "navigate-between-two-workspaces",
			steps: []steps{
				{workspace: "A", expected: "hello world\n", toWrite: "AA"},
				{workspace: "B", expected: "hello world\n", toWrite: "BB"},
				{workspace: "A", expected: "AA", toWrite: "AAaa"},
				{workspace: "B", expected: "BB", toWrite: "BBbb"},
				{workspace: "A", expected: "AAaa"},
				{workspace: "B", expected: "BBbb"},
				{workspace: "A", expected: "AAaa"},
				{workspace: "B", expected: "BBbb", toWriteNewFile: "stuff"},
				{workspace: "A", expected: "AAaa"},
				{workspace: "B", expected: "BBbb", expectInNewFile: "stuff"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			codebaseID := uuid.NewString()
			viewID := uuid.NewString()

			trunkPath := repoProvider.TrunkPath(codebaseID)
			viewPath := repoProvider.ViewPath(codebaseID, viewID)

			workspaceAID := uuid.New()
			workspaceBID := uuid.New()

			err := codebaseRepo.Create(codebase.Codebase{
				ID:              codebaseID,
				ShortCodebaseID: codebase.ShortCodebaseID(codebaseID), // not realistic
			})
			assert.NoError(t, err)
			assert.NoError(t, codebaseUserRepo.Create(codebase.CodebaseUser{
				ID:         uuid.NewString(),
				CodebaseID: codebaseID,
				UserID:     userID,
			}))
			assert.NoError(t, workspaceRepo.Create(workspaces.Workspace{
				ID:         workspaceAID.String(),
				CodebaseID: codebaseID,
				UserID:     userID,
			}))
			assert.NoError(t, workspaceRepo.Create(workspaces.Workspace{
				ID:         workspaceBID.String(),
				CodebaseID: codebaseID,
				UserID:     userID,
			}))
			err = viewRepo.Create(view.View{
				ID:          viewID,
				UserID:      userID,
				CodebaseID:  codebaseID,
				WorkspaceID: workspaceAID.String(),
			})
			assert.NoError(t, err)

			_, err = vcs.CreateBareRepoWithRootCommit(trunkPath)
			if err != nil {
				panic(err)
			}
			repoA, err := vcs.CloneRepo(trunkPath, viewPath)
			if err != nil {
				panic(err)
			}

			// Create common history
			assert.NoError(t, repoA.CheckoutBranchWithForce("sturdytrunk"))
			assert.NoError(t, ioutil.WriteFile(viewPath+"/a.txt", []byte("hello world\n"), 0666))
			_, err = repoA.AddAndCommit("Commit 1 (in A)")
			assert.NoError(t, err)
			assert.NoError(t, repoA.Push(zap.NewNop(), "sturdytrunk"))

			trunkLogEntries, err := repoA.LogBranch("sturdytrunk", 10)
			assert.NoError(t, err)

			// Create two branches
			assert.NoError(t, repoA.CreateNewBranchOnHEAD(workspaceAID.String()))
			assert.NoError(t, repoA.CreateNewBranchOnHEAD(workspaceBID.String()))

			assert.NoError(t, repoA.Push(zap.NewNop(), workspaceAID.String()))
			assert.NoError(t, repoA.Push(zap.NewNop(), workspaceBID.String()))

			for _, s := range tc.steps {
				var workspaceID string
				if s.workspace == "A" {
					workspaceID = workspaceAID.String()
				} else if s.workspace == "B" {
					workspaceID = workspaceBID.String()
				}

				_, err = viewResolver.OpenWorkspaceOnView(authCtx, resolvers.OpenViewArgs{
					Input: resolvers.OpenWorkspaceOnViewInput{
						WorkspaceID: graphql.ID(workspaceID),
						ViewID:      graphql.ID(viewID),
					},
				})
				assert.NoError(t, err)

				// Content as expected
				fileContent, err := ioutil.ReadFile(viewPath + "/a.txt")
				assert.NoError(t, err)
				assert.Equal(t, s.expected, string(fileContent))
				// New file content as expected
				if s.expectInNewFile != "" {
					fileContent, err := ioutil.ReadFile(viewPath + "/newfile.txt")
					assert.NoError(t, err)
					assert.Equal(t, s.expectInNewFile, string(fileContent))
				} else {
					// File not expected to be there
					assert.NoFileExists(t, viewPath+"/newfile.txt")
				}
				// No new commits
				wsLogEntries, err := repoA.LogBranch(workspaceID, 10)
				assert.NoError(t, err)
				assert.Equal(t, trunkLogEntries, wsLogEntries)
				// Write some new unsaved changes
				if s.toWrite != "" {
					assert.NoError(t, ioutil.WriteFile(viewPath+"/a.txt", []byte(s.toWrite), 0666))
				}
				if s.toWriteNewFile != "" {
					assert.NoError(t, ioutil.WriteFile(viewPath+"/newfile.txt", []byte(s.toWriteNewFile), 0666))
				}
			}

		})
	}
}

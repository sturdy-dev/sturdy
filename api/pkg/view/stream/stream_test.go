package stream_test

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"testing"

	"getsturdy.com/api/pkg/view/stream"
	"getsturdy.com/api/pkg/view/vcs"
	"getsturdy.com/api/vcs/provider"

	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/codebase"
	db_acl "getsturdy.com/api/pkg/codebase/acl/db"
	acl_provider "getsturdy.com/api/pkg/codebase/acl/provider"
	db_codebase "getsturdy.com/api/pkg/codebase/db"
	service_codebase "getsturdy.com/api/pkg/codebase/service"
	codebasevcs "getsturdy.com/api/pkg/codebase/vcs"
	"getsturdy.com/api/pkg/db"
	"getsturdy.com/api/pkg/events"
	"getsturdy.com/api/pkg/internal/sturdytest"
	db_suggestion "getsturdy.com/api/pkg/suggestions/db"
	service_suggestion "getsturdy.com/api/pkg/suggestions/service"
	"getsturdy.com/api/pkg/users"
	db_user "getsturdy.com/api/pkg/users/db"
	service_user "getsturdy.com/api/pkg/users/service"
	"getsturdy.com/api/pkg/view"
	db_view "getsturdy.com/api/pkg/view/db"
	"getsturdy.com/api/pkg/workspaces"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"
	workspacevcs "getsturdy.com/api/pkg/workspaces/vcs"
	"getsturdy.com/api/vcs/executor"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// TestStream is a fairly basic test that tests the streaming functionality
// As of now, it only tests that messages are received, and does not test their contents
func TestStream(t *testing.T) {
	if os.Getenv("E2E_TEST") == "" {
		t.SkipNow()
	}

	testCases := []struct {
		name     string
		withView bool
	}{
		{
			name:     "with view",
			withView: true,
		},
		{
			name:     "without view",
			withView: false,
		},
	}

	repoProvider := newRepoProvider(t)
	logger, _ := zap.NewDevelopment()
	executorProvider := executor.NewProvider(logger, repoProvider)
	viewUpdates := events.NewInMemory()

	d, err := db.Setup(
		sturdytest.PsqlDbSourceForTesting(),
	)
	if err != nil {
		panic(err)
	}

	viewRepo := db_view.NewRepo(d)
	workspaceRepo := db_workspaces.NewRepo(d)
	// snapshotRepo := db_snapshots.NewRepo(d)
	userRepo := db_user.NewRepo(d)
	codebaseRepo := db_codebase.NewRepo(d)
	aclRepo := db_acl.NewACLRepository(d)
	codebaseUserRepo := db_codebase.NewCodebaseUserRepo(d)
	aclProvider := acl_provider.New(aclRepo, codebaseUserRepo, userRepo)
	userService := service_user.New(zap.NewNop(), userRepo, nil)
	suggestionsDB := db_suggestion.New(d)

	workspaceService := service_workspace.New(
		logger,
		nil,
		workspaceRepo,
		workspaceRepo,

		userRepo,
		nil,

		nil,
		nil,

		nil,
		executorProvider,
		nil,
		nil,
		nil,
		nil,
	)

	codebaseService := service_codebase.New(codebaseRepo, codebaseUserRepo, workspaceService, nil, logger, executorProvider, nil, nil)
	authService := service_auth.New(codebaseService, userService, nil, aclProvider, nil /*organizationService*/)

	suggestionsService := service_suggestion.New(
		logger,
		suggestionsDB,
		workspaceService,
		executorProvider,
		nil,
		nil,
		nil,
		nil,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			codebaseID := uuid.NewString()
			err := codebasevcs.Create(repoProvider, codebaseID)
			assert.NoError(t, err)

			workspaceID := uuid.NewString()
			trunkRepo, err := repoProvider.TrunkRepo(codebaseID)
			assert.NoError(t, err)
			err = workspacevcs.Create(trunkRepo, workspaceID)
			assert.NoError(t, err)

			viewID := uuid.NewString()
			err = vcs.Create(repoProvider, codebaseID, workspaceID, viewID)
			assert.NoError(t, err)

			userID := users.ID(uuid.NewString())

			err = userRepo.Create(&users.User{ID: userID, Email: userID.String() + "@test.getsturdy.com"})
			assert.NoError(t, err)

			err = codebaseRepo.Create(codebase.Codebase{
				ID:              codebaseID,
				ShortCodebaseID: codebase.ShortCodebaseID(codebaseID), // this is not realistic, but it works
			})
			assert.NoError(t, err)

			ws := workspaces.Workspace{
				ID:         workspaceID,
				CodebaseID: codebaseID,
				ViewID:     &viewID,
				UserID:     userID,
			}
			err = workspaceRepo.Create(ws)
			assert.NoError(t, err)

			var vw *view.View

			if tc.withView {
				vw = &view.View{
					ID:          viewID,
					WorkspaceID: workspaceID,
					CodebaseID:  codebaseID,
					UserID:      userID,
				}
				err = viewRepo.Create(*vw)
				assert.NoError(t, err)
			}

			var wg sync.WaitGroup

			// If the client has disconnected
			done := make(chan bool)

			wg.Add(1)
			go func() {
				events, err := stream.Stream(
					context.Background(),
					logger,
					&ws,
					vw,
					done,
					viewUpdates,
					authService,
					workspaceService,
					suggestionsService,
				)
				if !assert.NoError(t, err) {
					log.Println(err)
					return
				}

				// Drain events
				for event := range events {
					log.Println(event)
					t.Logf("event=%+v", event)
				}

				wg.Done()
			}()

			// Disconnect after some time
			wg.Add(1)
			go func() {
				done <- true
				wg.Done()
			}()

			wg.Wait()
		})
	}
}

func newRepoProvider(t *testing.T) provider.RepoProvider {
	reposBasePath, err := ioutil.TempDir(os.TempDir(), "mash")
	assert.NoError(t, err)
	return provider.New(reposBasePath, "")
}

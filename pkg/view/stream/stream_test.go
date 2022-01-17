package stream

import (
	"context"
	"log"
	"os"
	"sync"
	"testing"

	"mash/pkg/analytics/disabled"
	"mash/pkg/view/vcs"

	"mash/db"
	service_auth "mash/pkg/auth/service"
	"mash/pkg/codebase"
	db_acl "mash/pkg/codebase/acl/db"
	acl_provider "mash/pkg/codebase/acl/provider"
	db_codebase "mash/pkg/codebase/db"
	service_codebase "mash/pkg/codebase/service"
	codebasevcs "mash/pkg/codebase/vcs"
	"mash/pkg/internal/sturdytest"
	db_suggestion "mash/pkg/suggestions/db"
	service_suggestion "mash/pkg/suggestions/service"
	"mash/pkg/user"
	db_user "mash/pkg/user/db"
	service_user "mash/pkg/user/service"
	"mash/pkg/view"
	db_view "mash/pkg/view/db"
	"mash/pkg/view/events"
	"mash/pkg/workspace"
	db_workspace "mash/pkg/workspace/db"
	service_workspace "mash/pkg/workspace/service"
	workspacevcs "mash/pkg/workspace/vcs"
	"mash/vcs/executor"

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
		true,
		"file://../../../db/migrations",
	)
	if err != nil {
		panic(err)
	}

	viewRepo := db_view.NewRepo(d)
	workspaceRepo := db_workspace.NewRepo(d)
	// snapshotRepo := db_snapshots.NewRepo(d)
	userRepo := db_user.NewRepo(d)
	codebaseRepo := db_codebase.NewRepo(d)
	aclRepo := db_acl.NewACLRepository(d)
	codebaseUserRepo := db_codebase.NewCodebaseUserRepo(d)
	aclProvider := acl_provider.New(aclRepo, codebaseUserRepo, userRepo)
	codebaseService := service_codebase.New(codebaseRepo, codebaseUserRepo)
	userService := service_user.New(zap.NewNop(), userRepo, nil /*jwtService*/, nil /*onetime*/, nil /*emailsender*/, disabled.NewClient())
	authService := service_auth.New(codebaseService, userService, nil, aclProvider, nil /*organizationService*/)
	suggestionsDB := db_suggestion.New(d)

	workspaceService := service_workspace.New(
		logger,
		disabled.NewClient(),

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

			userID := uuid.NewString()

			err = userRepo.Create(&user.User{ID: userID, Email: userID + "@test.getsturdy.com"})
			assert.NoError(t, err)

			err = codebaseRepo.Create(codebase.Codebase{
				ID:              codebaseID,
				ShortCodebaseID: codebase.ShortCodebaseID(codebaseID), // this is not realistic, but it works
			})
			assert.NoError(t, err)

			ws := workspace.Workspace{
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
				events, err := Stream(
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

package db_test

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"mash/pkg/db"
	"mash/pkg/internal/sturdytest"
	"mash/pkg/workspace/watchers"
	watchers_db "mash/pkg/workspace/watchers/db"

	"github.com/stretchr/testify/assert"
)

var implementations = []func() watchers_db.Repository{
	func() watchers_db.Repository {
		return watchers_db.NewInMemory()
	},
}

var tests = []func(*testing.T, watchers_db.Repository){
	ShouldListIfWatching,
	ShouldNotListIfIgnored,
	ShouldReturnByUserIDWorkspaceID,
}

func ShouldListIfWatching(t *testing.T, repo watchers_db.Repository) {
	workspaceID := "workspace-id"
	assert.NoError(t, repo.Create(context.Background(), &watchers.Watcher{
		UserID:      "user1",
		WorkspaceID: workspaceID,
		CreatedAt:   time.Now(),
		Status:      watchers.StatusWatching,
	}))

	assert.NoError(t, repo.Create(context.Background(), &watchers.Watcher{
		UserID:      "user2",
		WorkspaceID: workspaceID,
		CreatedAt:   time.Now(),
		Status:      watchers.StatusIgnored,
	}))

	ww, err := repo.ListWatchingByWorkspaceID(context.Background(), workspaceID)
	assert.NoError(t, err)
	if assert.Len(t, ww, 1) {
		assert.Equal(t, "user1", ww[0].UserID)
	}
}

func ShouldNotListIfIgnored(t *testing.T, repo watchers_db.Repository) {
	workspaceID := "workspace-id"
	assert.NoError(t, repo.Create(context.Background(), &watchers.Watcher{
		UserID:      "user1",
		WorkspaceID: workspaceID,
		CreatedAt:   time.Now(),
		Status:      watchers.StatusWatching,
	}))

	assert.NoError(t, repo.Create(context.Background(), &watchers.Watcher{
		UserID:      "user1",
		WorkspaceID: workspaceID,
		CreatedAt:   time.Now(),
		Status:      watchers.StatusIgnored,
	}))

	assert.NoError(t, repo.Create(context.Background(), &watchers.Watcher{
		UserID:      "user2",
		WorkspaceID: workspaceID,
		CreatedAt:   time.Now(),
		Status:      watchers.StatusIgnored,
	}))

	ww, err := repo.ListWatchingByWorkspaceID(context.Background(), workspaceID)
	assert.NoError(t, err)
	assert.Len(t, ww, 0)
}

func ShouldReturnByUserIDWorkspaceID(t *testing.T, repo watchers_db.Repository) {
	workspaceID := "workspace-id"
	assert.NoError(t, repo.Create(context.Background(), &watchers.Watcher{
		UserID:      "user1",
		WorkspaceID: workspaceID,
		CreatedAt:   time.Now(),
		Status:      watchers.StatusWatching,
	}))

	assert.NoError(t, repo.Create(context.Background(), &watchers.Watcher{
		UserID:      "user2",
		WorkspaceID: workspaceID,
		CreatedAt:   time.Now(),
		Status:      watchers.StatusIgnored,
	}))

	w, err := repo.GetByUserIDAndWorkspaceID(context.Background(), "user1", workspaceID)
	if assert.NoError(t, err) {
		assert.Equal(t, "user1", w.UserID)
		assert.Equal(t, workspaceID, w.WorkspaceID)
		assert.Equal(t, watchers.StatusWatching, w.Status)
	}

	w, err = repo.GetByUserIDAndWorkspaceID(context.Background(), "user2", workspaceID)
	if assert.NoError(t, err) {
		assert.Equal(t, "user2", w.UserID)
		assert.Equal(t, workspaceID, w.WorkspaceID)
		assert.Equal(t, watchers.StatusIgnored, w.Status)
	}

	w, err = repo.GetByUserIDAndWorkspaceID(context.Background(), "user3", workspaceID)
	if assert.ErrorIs(t, err, sql.ErrNoRows) {
		assert.Nil(t, w)
	}

	w, err = repo.GetByUserIDAndWorkspaceID(context.Background(), "user1", "workspace-id-2")
	if assert.ErrorIs(t, err, sql.ErrNoRows) {
		assert.Nil(t, w)
	}
}

func TestMain(m *testing.M) {
	defer m.Run()

	if os.Getenv("E2E_TEST") == "" {
		return
	}

	// register real db implementation
	sqldb, err := db.Setup(sturdytest.PsqlDbSourceForTesting())
	if err != nil {
		panic(err)
	}
	databaseImplementation := func() watchers_db.Repository { return watchers_db.NewDB(sqldb) }

	implementations = append(implementations, databaseImplementation)
}

// runs all tests for a all implementations
func TestImplementations(t *testing.T) {
	for _, test := range tests {
		t.Run(funcName(test), func(t *testing.T) {
			for _, repoProvider := range implementations {
				repo := repoProvider()
				t.Run(implName(repo), func(t *testing.T) {
					test(t, repo)
				})
			}
		})
	}
}

func funcName(v interface{}) string {
	pc := reflect.ValueOf(v).Pointer()
	nameFull := runtime.FuncForPC(pc).Name()
	nameEnd := filepath.Ext(nameFull)
	name := strings.TrimPrefix(nameEnd, ".")
	return name
}

func implName(v interface{}) string {
	nameFull := reflect.TypeOf(v).String()
	nameEnd := filepath.Ext(nameFull)
	name := strings.TrimPrefix(nameEnd, ".")
	return name
}

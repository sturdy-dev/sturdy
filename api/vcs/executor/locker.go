package executor

import (
	"fmt"
	"path/filepath"
	"sync"

	"getsturdy.com/api/vcs/provider"

	"github.com/gofrs/flock"
)

type locker struct {
	provider provider.RepoProvider

	locksGuard *sync.Mutex
	locks      map[string]lock
}

type lock interface {
	Lock() error
	Unlock() error
	RLock() error
	RUnlock() error
}

type mutexLock struct {
	mu *sync.RWMutex
}

func (ml *mutexLock) Lock() error {
	ml.mu.Lock()
	return nil
}

func (ml *mutexLock) Unlock() error {
	ml.mu.Unlock()
	return nil
}

func (ml *mutexLock) RLock() error {
	ml.mu.RLock()
	return nil
}

func (ml *mutexLock) RUnlock() error {
	ml.mu.RUnlock()
	return nil
}

func newLocker(provider provider.RepoProvider) *locker {
	return &locker{
		provider: provider,

		locksGuard: &sync.Mutex{},
		locks:      map[string]lock{},
	}
}

// use .git directory to store lock file in, because it's the only place that is not synced with the users' end.
const lockFileName = ".git/sturdy.lock"

// Returns a mutex for the given codebase and view.
func (l *locker) Get(codebaseID string, viewID *string) lock {
	key := l.key(codebaseID, viewID)

	l.locksGuard.Lock()
	defer l.locksGuard.Unlock()

	if m, ok := l.locks[key]; ok {
		return m
	}

	// codebases are locked using in-memory mutexes
	if viewID == nil {
		mutex := &mutexLock{&sync.RWMutex{}}
		l.locks[key] = mutex
		return mutex
	}

	// for views, we use file locks to synchronize with mutagen process
	lockFile := filepath.Join(l.provider.ViewPath(codebaseID, *viewID), lockFileName)
	lock := &FileLock{lock: flock.New(lockFile)}
	l.locks[key] = lock
	return lock
}

func (l *locker) key(codebaseID string, viewID *string) string {
	if viewID == nil {
		return fmt.Sprintf("%s/trunk", codebaseID)
	}
	return fmt.Sprintf("%s/%s", codebaseID, *viewID)
}

func (l *locker) GetInMemory(codebaseID string, viewID *string) lock {
	key := l.key(codebaseID, viewID) + "-inmemory"

	l.locksGuard.Lock()
	defer l.locksGuard.Unlock()

	if m, ok := l.locks[key]; ok {
		return m
	}

	mutex := &mutexLock{&sync.RWMutex{}}
	l.locks[key] = mutex
	return mutex
}

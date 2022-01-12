// IMPORTANT: do not forget to update ./testdata/ files when making changes here

package executor

import (
	"errors"
	"os"
	"sync"

	"github.com/gofrs/flock"
)

type FileLock struct {
	mu      sync.RWMutex // lock within this process
	lock    *flock.Flock // lock between processes
	countMx sync.Mutex

	count uint
}

func New(filename string) *FileLock {
	return &FileLock{
		lock: flock.New(filename),
	}
}

func (fl *FileLock) Lock() error {
	fl.mu.Lock()
	if err := fl.lock.Lock(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	return nil
}

func (fl *FileLock) Unlock() error {
	defer fl.mu.Unlock()
	if err := fl.lock.Unlock(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	return nil
}

func (fl *FileLock) RLock() error {
	fl.countMx.Lock()
	fl.count++
	fl.countMx.Unlock()

	fl.mu.RLock()
	if err := fl.lock.RLock(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	return nil
}

func (fl *FileLock) RUnlock() error {
	fl.countMx.Lock()
	fl.count--
	if fl.count == 0 {
		if err := fl.lock.Unlock(); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				fl.mu.RUnlock()
				fl.countMx.Unlock()
				return nil
			}
			fl.mu.RUnlock()
			fl.countMx.Unlock()
			return err

		}
	}
	fl.mu.RUnlock()
	fl.countMx.Unlock()
	return nil
}

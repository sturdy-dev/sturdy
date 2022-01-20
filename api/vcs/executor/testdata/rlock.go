package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gofrs/flock"
)

var (
	filepath = flag.String("filepath", "", "")
	timeout  = flag.Duration("timeout", 1, "")
	hold     = flag.Duration("hold", 1, "")
	prefix   = flag.String("prefix", "lock.go", "")
)

func main() {
	flag.Parse()

	if *filepath == "" {
		log.Fatal(*prefix, "filepath is required")
	}

	if *timeout == 1 {
		log.Fatal(*prefix, "timeout is required")
	}

	if *hold == 1 {
		log.Fatal(*prefix, "hold is required")
	}

	lock := NewFileLock(*filepath)

	lockAcquired := make(chan struct{})
	go func() {
		log.Println(*prefix, "acquiring rlock")
		lock.RLock()
		close(lockAcquired)
	}()

	select {
	case <-lockAcquired:
		log.Println(*prefix, "lock acquired")
		<-time.After(*hold)
		log.Println(*prefix, "releasing lock")
		lock.Unlock()
		log.Println(*prefix, "lock released")
	case <-time.After(*timeout):
		log.Fatal(*prefix, " timeout")
	}
}

// this is a copy of lock.go to simplify testing
// do not forget to update it if needed

type FileLock struct {
	mu      sync.RWMutex // lock within this process
	lock    *flock.Flock // lock between processes
	countMx sync.Mutex

	count uint
}

func NewFileLock(filename string) *FileLock {
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

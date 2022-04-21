package executor

import (
	"path/filepath"
	"testing"
	"time"
)

func TestInProcess_RLockNonBlocksRLock(t *testing.T) {
	dir := t.TempDir()
	lock := New(filepath.Join(dir, ".lock"))

	lockSecond := make(chan struct{})
	secondLocked := make(chan struct{})
	go func() {
		lock.RLock()
		close(lockSecond)
	}()
	go func() {
		<-lockSecond
		lock.RLock()
		close(secondLocked)
	}()
	select {
	case <-time.Tick(10 * time.Millisecond):
		t.Fatal("rlock is blocking rlock")
	case <-secondLocked:
		// success, unlock unblocks lock
	}
}

func TestInProcess_UnlockUnblocksLock(t *testing.T) {
	dir := t.TempDir()
	lock := New(filepath.Join(dir, ".lock"))

	lockSecond := make(chan struct{})
	secondLocked := make(chan struct{})
	unlockFirst := make(chan struct{})
	go func() {
		lock.Lock()
		close(lockSecond)
		<-unlockFirst
		lock.Unlock()
	}()
	go func() {
		<-lockSecond
		close(unlockFirst)
		lock.Lock()
		close(secondLocked)
	}()
	select {
	case <-time.Tick(10 * time.Millisecond):
		t.Fatal("unlock is not unblocking lock")
	case <-secondLocked:
		// success, unlock unblocks lock
	}
}

func TestInProcess_LockBlocksLock(t *testing.T) {
	dir := t.TempDir()
	lock := New(filepath.Join(dir, ".lock"))
	lockSecond := make(chan struct{})
	secondLocked := make(chan struct{})
	go func() {
		lock.Lock()
		close(lockSecond)
	}()
	go func() {
		<-lockSecond
		lock.Lock()
		close(secondLocked)
	}()
	select {
	case <-time.Tick(10 * time.Millisecond):
		// success, lock is blocking lock
	case <-secondLocked:
		t.Fatal("lock is not blocking lock")
	}
}

func TestInProcess_RLockBlocksLock(t *testing.T) {
	dir := t.TempDir()
	lock := New(filepath.Join(dir, ".lock"))

	lockSecond := make(chan struct{})
	secondLocked := make(chan struct{})
	go func() {
		lock.RLock()
		close(lockSecond)
	}()
	go func() {
		<-lockSecond
		lock.Lock()
		close(secondLocked)
	}()
	select {
	case <-time.Tick(10 * time.Millisecond):
		// success, rlock is blocking lock
	case <-secondLocked:
		t.Fatal("rlock is not blocking lock")
	}
}
func TestInProcess_RUnlockUnblocksLock(t *testing.T) {
	dir := t.TempDir()

	lock := New(filepath.Join(dir, ".lock"))

	lockSecond := make(chan struct{})
	secondLocked := make(chan struct{})
	runlockFirst := make(chan struct{})
	go func() {
		lock.RLock()
		close(lockSecond)
		<-runlockFirst
		lock.RUnlock()
	}()
	go func() {
		<-lockSecond
		close(runlockFirst)
		lock.Lock()
		close(secondLocked)
	}()
	select {
	case <-secondLocked:
		// success, runlock is unblocking lock
	case <-time.Tick(10 * time.Millisecond):
		t.Fatal("rlock is not unblocking lock")
	}
}

func TestInProcess_LockBlocksRLock(t *testing.T) {
	dir := t.TempDir()
	lock := New(filepath.Join(dir, ".lock"))
	rlockSecond := make(chan struct{})
	secondRLocked := make(chan struct{})
	go func() {
		lock.Lock()
		close(rlockSecond)
	}()
	go func() {
		<-rlockSecond
		lock.RLock()
		close(secondRLocked)
	}()
	select {
	case <-time.Tick(10 * time.Millisecond):
		// success, lock is blocking rlock
	case <-secondRLocked:
		t.Fatal("lock is not blocking rlock")
	}
}

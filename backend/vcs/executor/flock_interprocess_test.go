package executor

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type lockType string

var (
	writeLock = lockType("lock")
	readLock  = lockType("rlock")
)

func notBlocks(t *testing.T, l1, l2 lockType) {
	lockOne := exec.Command(
		"go", "run", "./testdata/lock.go",
		"--filepath", ".lock",
		"--timeout", "10ms",
		"--hold", "1s",
		"--prefix", "one",
		"--mode", string(l1),
	)
	lockOne.Stdout = os.Stdout
	lockOne.Stderr = os.Stderr

	lockTwo := exec.Command(
		"go", "run", "./testdata/lock.go",
		"--filepath", ".lock",
		"--timeout", "10ms",
		"--hold", "1s",
		"--prefix", "two",
		"--mode", string(l2),
	)
	lockTwo.Stdout = os.Stdout
	lockTwo.Stderr = os.Stderr

	assert.NoError(t, lockOne.Start())
	time.Sleep(time.Millisecond * 100)
	assert.NoError(t, lockTwo.Start())

	assert.NoError(t, lockTwo.Wait())
	assert.NoError(t, lockOne.Wait())
}

func blocks(t *testing.T, l1, l2 lockType) {
	lockOne := exec.Command(
		"go", "run", "./testdata/lock.go",
		"--filepath", ".lock",
		"--timeout", "1s",
		"--hold", "1s",
		"--prefix", "one",
		"--mode", string(l1),
	)
	lockOne.Stdout = os.Stdout
	lockOne.Stderr = os.Stderr

	lockTwo := exec.Command(
		"go", "run", "./testdata/lock.go",
		"--filepath", ".lock",
		"--timeout", "10ms",
		"--hold", "1s",
		"--prefix", "two",
		"--mode", string(l2),
	)
	lockTwo.Stdout = os.Stdout
	lockTwo.Stderr = os.Stderr

	assert.NoError(t, lockOne.Start())
	time.Sleep(time.Millisecond * 100)
	assert.NoError(t, lockTwo.Start())

	assert.Error(t, lockTwo.Wait())
	assert.NoError(t, lockOne.Wait())
}

func unblocks(t *testing.T, l1, l2 lockType) {
	lockOne := exec.Command(
		"go", "run", "./testdata/lock.go",
		"--filepath", ".lock",
		"--timeout", "1s",
		"--hold", "500ms",
		"--prefix", "one",
		"--mode", string(l1),
	)
	lockOne.Stdout = os.Stdout
	lockOne.Stderr = os.Stderr

	lockTwo := exec.Command(
		"go", "run", "./testdata/lock.go",
		"--filepath", ".lock",
		"--timeout", "1s",
		"--hold", "1ms",
		"--prefix", "two",
		"--mode", string(l2),
	)
	lockTwo.Stdout = os.Stdout
	lockTwo.Stderr = os.Stderr

	assert.NoError(t, lockOne.Start())
	time.Sleep(time.Millisecond * 100)
	assert.NoError(t, lockTwo.Start())

	assert.NoError(t, lockOne.Wait())
	assert.NoError(t, lockTwo.Wait())
}

func TestInterProcess_RLockNonBlocksRLock(t *testing.T) {
	notBlocks(t, readLock, readLock)
}

func TestInterProcess_UnlockUnblocksLock(t *testing.T) {
	unblocks(t, writeLock, writeLock)
}

func TestInterProcess_LockBlocksLock(t *testing.T) {
	blocks(t, writeLock, writeLock)
}

func TestInterProcess_RLockBlocksLock(t *testing.T) {
	blocks(t, readLock, writeLock)
}

func TestInterProcess_RUnlockUnblocksLock(t *testing.T) {
	unblocks(t, readLock, writeLock)
}

func TestInterProcess_LockBlocksRLock(t *testing.T) {
	blocks(t, writeLock, readLock)
}

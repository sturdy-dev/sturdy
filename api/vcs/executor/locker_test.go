package executor

import (
	"testing"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/vcs/testutil"

	"github.com/stretchr/testify/assert"
)

func TestLocker_Get__returns_same_lock(t *testing.T) {
	repoProvider := testutil.TestingRepoProvider(t)
	l := newLocker(repoProvider)

	codebaseID := codebases.ID("cb1")
	var viewID *string

	lock1 := l.Get(codebaseID, viewID)
	lock2 := l.Get(codebaseID, viewID)
	assert.Equal(t, lock1, lock2)
}

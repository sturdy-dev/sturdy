package inmemory

import (
	"getsturdy.com/api/pkg/di"
	db_suggestions "getsturdy.com/api/pkg/suggestions/db"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
)

func TestModul(c *di.Container) {
	c.Register(NewInMemoryAclRepo)
	c.Register(NewInMemorySnapshotRepo)
	c.Register(NewInMemoryViewRepo)
	c.Register(db_workspaces.NewMemory)
	c.Register(db_suggestions.NewMemory)
}

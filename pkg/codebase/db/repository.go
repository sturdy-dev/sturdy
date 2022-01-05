package db

import "mash/pkg/codebase"

type CodebaseRepository interface {
	Create(entity codebase.Codebase) error
	Get(id string) (*codebase.Codebase, error)
	GetAllowArchived(id string) (*codebase.Codebase, error)
	GetByInviteCode(inviteCode string) (*codebase.Codebase, error)
	GetByShortID(shortID string) (*codebase.Codebase, error)
	Update(entity *codebase.Codebase) error
}

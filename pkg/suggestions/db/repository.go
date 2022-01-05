package db

import (
	"context"

	"mash/pkg/suggestions"
)

type Repository interface {
	Create(context.Context, *suggestions.Suggestion) error
	Update(context.Context, *suggestions.Suggestion) error
	GetByID(context.Context, suggestions.ID) (*suggestions.Suggestion, error)
	GetByWorkspaceID(context.Context, string) (*suggestions.Suggestion, error)
	ListForWorkspaceID(context.Context, string) ([]*suggestions.Suggestion, error)
	ListBySnapshotID(context.Context, string) ([]*suggestions.Suggestion, error)
}

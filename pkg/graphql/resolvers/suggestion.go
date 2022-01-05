package resolvers

import (
	"context"

	"mash/pkg/suggestions"

	"github.com/graph-gophers/graphql-go"
)

type SuggestionRootResolver interface {
	// Internal
	InternalSuggestion(context.Context, *suggestions.Suggestion) (SuggestionResolver, error)
	InternalSuggestionByID(context.Context, suggestions.ID) (SuggestionResolver, error)

	// Mutations
	CreateSuggestion(context.Context, CreateSuggestionArgs) (SuggestionResolver, error)
	DismissSuggestion(context.Context, DismissSuggestionArgs) (SuggestionResolver, error)
	ApplySuggestionHunks(context.Context, ApplySuggestionHunksArgs) (SuggestionResolver, error)
	DismissSuggestionHunks(context.Context, DismissSuggestionHunksArgs) (SuggestionResolver, error)

	// Subscriptions
	UpdatedSuggestion(context.Context, UpdatedSuggestionArgs) (chan SuggestionResolver, error)
}

type ApplySuggestionHunksArgs struct {
	Input ApplySuggestionHunksInput
}

type ApplySuggestionHunksInput struct {
	ID      graphql.ID
	HunkIDs []string
}

type DismissSuggestionHunksArgs struct {
	Input DismissSuggestionHunksInput
}

type DismissSuggestionHunksInput struct {
	ID      graphql.ID
	HunkIDs []string
}

type CreateSuggestionArgs struct {
	Input CreateSuggestionInput
}

type CreateSuggestionInput struct {
	WorkspaceID graphql.ID
}

type UpdatedSuggestionArgs struct {
	WorkspaceID graphql.ID
}

type DismissSuggestionArgs struct {
	Input DismissSuggestionInput
}

type DismissSuggestionInput struct {
	ID graphql.ID
}

type SuggestionResolver interface {
	ID() graphql.ID
	Author(context.Context) (AuthorResolver, error)
	Workspace(context.Context) (WorkspaceResolver, error)
	For(context.Context) (WorkspaceResolver, error)
	Diffs(context.Context) ([]FileDiffResolver, error)
	CreatedAt() int32
	DismissedAt() *int32
}

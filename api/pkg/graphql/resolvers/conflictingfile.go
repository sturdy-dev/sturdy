package resolvers

import (
	"github.com/graph-gophers/graphql-go"
)

type ConflictingFileResolver interface {
	ID() graphql.ID
	Path() string
	WorkspaceDiff() (FileDiffResolver, error)
	TrunkDiff() (FileDiffResolver, error)
}

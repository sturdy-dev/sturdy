package graphql

import (
	"context"
	"errors"
	"github.com/graph-gophers/graphql-go"
	"log"
	gqlerrors "mash/pkg/graphql/errors"
	"mash/pkg/graphql/resolvers"
	"path"
	"sort"
	"strings"
)

type directoryResolver struct {
	codebaseID   string
	path         string
	children     []string
	rootResolver resolvers.FileRootResolver
}

func (r *directoryResolver) ToFile() (resolvers.FileResolver, bool) {
	return nil, false
}

func (r *directoryResolver) ToDirectory() (resolvers.DirectoryResolver, bool) {
	return r, true
}

func (r *directoryResolver) ID() graphql.ID {
	return graphql.ID(r.codebaseID + "-" + r.path)
}

func (r *directoryResolver) Path() string {
	return r.path
}

func (r *directoryResolver) Children(ctx context.Context) []resolvers.FileOrDirectoryResolver {
	var children []resolvers.FileOrDirectoryResolver

	for _, child := range r.children {
		file, err := r.rootResolver.InternalFile(ctx, r.codebaseID, child)
		if err != nil {
			log.Println("failed to open child", err)
		} else {
			children = append(children, file)
		}
	}

	sort.Slice(children, func(i, j int) bool {
		a, b := children[i], children[j]

		_, aIsFile := a.ToFile()
		_, bIsFile := b.ToFile()

		switch {
		case !aIsFile && bIsFile:
			return true

		case !bIsFile && aIsFile:
			return false

		default:
			return strings.ToLower(a.Path()) < strings.ToLower(b.Path())
		}
	})

	return children
}

func (r *directoryResolver) Readme(ctx context.Context) (resolvers.FileResolver, error) {
	// GitHub supported names:
	// https://github.com/github/markup/blob/master/README.md
	fileResolver, err := r.rootResolver.InternalFile(
		ctx,
		r.codebaseID,
		path.Join(r.path, "README.md"),
		path.Join(r.path, "README.mkdn"),
		path.Join(r.path, "README.mdown"),
		path.Join(r.path, "README.markdown"),
	)
	switch {
	case err == nil:
		if file, ok := fileResolver.ToFile(); ok {
			return file, nil
		} else {
			return nil, nil
		}
	case errors.Is(err, gqlerrors.ErrNotFound):
		return nil, nil
	default:
		return nil, gqlerrors.Error(err)
	}
}

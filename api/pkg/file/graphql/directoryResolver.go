package graphql

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path"
	"sort"
	"strings"

	"github.com/graph-gophers/graphql-go"

	"getsturdy.com/api/pkg/codebases"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
)

type directoryResolver struct {
	codebase *codebases.Codebase
	path     string
	children []string

	rootResolver *fileRootResolver
}

func (r *directoryResolver) ToFile() (resolvers.FileResolver, bool) {
	return nil, false
}

func (r *directoryResolver) ToDirectory() (resolvers.DirectoryResolver, bool) {
	return r, true
}

func (r *directoryResolver) ID() graphql.ID {
	return graphql.ID(fmt.Sprintf("%s-%s", r.codebase.ID, r.path))
}

func (r *directoryResolver) Path() string {
	return r.path
}

func (r *directoryResolver) Children(ctx context.Context) ([]resolvers.FileOrDirectoryResolver, error) {
	var children []resolvers.FileOrDirectoryResolver

	allower, err := r.rootResolver.authService.GetAllower(ctx, r.codebase)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	for _, child := range r.children {
		if !allower.IsAllowed(child, false) {
			continue
		}

		file, err := r.rootResolver.InternalFile(ctx, r.codebase, child)
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

	return children, nil
}

func (r *directoryResolver) Readme(ctx context.Context) (resolvers.FileResolver, error) {
	// GitHub supported names:
	// https://github.com/github/markup/blob/master/README.md
	fileResolver, err := r.rootResolver.InternalFile(
		ctx,
		r.codebase,
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

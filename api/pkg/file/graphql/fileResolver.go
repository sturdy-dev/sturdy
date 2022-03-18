package graphql

import (
	"fmt"
	"path"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/graphql/resolvers"

	"github.com/graph-gophers/graphql-go"
)

type fileResolver struct {
	codebaseID codebases.ID
	path       string
	contents   []byte
}

func (r *fileResolver) ToFile() (resolvers.FileResolver, bool) {
	return r, true
}

func (r *fileResolver) ToDirectory() (resolvers.DirectoryResolver, bool) {
	return nil, false
}

func (r *fileResolver) ID() graphql.ID {
	return graphql.ID(fmt.Sprintf("%s-%s", r.codebaseID, r.path))
}

func (r *fileResolver) Path() string {
	return r.path
}

func (r *fileResolver) Contents() string {
	return string(r.contents)
}

func (r *fileResolver) MimeType() string {
	switch path.Ext(r.path) {
	case ".md", ".markdown", ".mdown":
		return "text/markdown"
	case ".json":
		return "application/json"
	case ".yaml":
		return "application/x-yaml"
	case ".js":
		return "application/javascript"
	case ".go":
		return "application/x-go"
	case ".py":
		return "application/x-py"
	case ".txt":
		return "text/plain"
	case ".html":
		return "text/html"
	case ".xml":
		return "text/xml"
	default:
		return "application/octet-stream"
	}
}

package uploader

import (
	"context"
	"fmt"
	"io"
	"path"

	"getsturdy.com/api/pkg/blobs"
	service_blobs "getsturdy.com/api/pkg/blobs/service"
	"getsturdy.com/api/pkg/users/avatars"
)

type Blobs struct {
	blobsService *service_blobs.Service

	baseURL   string
	urlPrefix string
}

func NewBlobs(cfg *Configuration, blobsService *service_blobs.Service) (*Blobs, error) {
	scheme := cfg.URL.Scheme
	host := cfg.URL.Host

	blobs := &Blobs{
		blobsService: blobsService,
		urlPrefix:    path.Join(cfg.URL.Path, "/"),
	}

	if (scheme == "" || host == "") && scheme+host != "" {
		return nil, fmt.Errorf("invalid URL: must provide either a full URL, or an absolute path")
	} else if scheme != "" {
		blobs.baseURL = fmt.Sprintf("%s://%s", scheme, host)
	}

	return blobs, nil
}

func (p *Blobs) Upload(ctx context.Context, key string, file io.Reader) (*avatars.Avatar, error) {
	if err := p.blobsService.Store(ctx, blobs.ID(key), file); err != nil {
		return nil, fmt.Errorf("failed to store blob: %w", err)
	}
	return &avatars.Avatar{
		URL: p.baseURL + path.Join(p.urlPrefix, "/v3/blobs/", key),
	}, nil
}

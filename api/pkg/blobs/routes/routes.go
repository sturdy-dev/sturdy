package routes

import (
	"errors"
	"net/http"
	"path"

	"getsturdy.com/api/pkg/blobs"
	service_blob "getsturdy.com/api/pkg/blobs/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Register(rg *gin.RouterGroup, logger *zap.Logger, blobService *service_blob.Service) {
	logger = logger.With(zap.String("handler", "routes/blobs"))
	rg.GET("/:key", gin.WrapF(Get(logger, blobService)))
}

func Get(logger *zap.Logger, blobService *service_blob.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := path.Base(r.URL.Path)

		blob, err := blobService.Fetch(r.Context(), blobs.ID(key))
		if errors.Is(err, service_blob.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			w.Header().Set("Content-Type", http.DetectContentType(blob.Data))
			if _, err := w.Write(blob.Data); err != nil {
				logger.Error("failed to write blob", zap.Error(err))
			}
		}
	}
}

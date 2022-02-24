package routes

import (
	"bytes"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/img"

	"github.com/google/uuid"

	"getsturdy.com/api/pkg/users/avatars/uploader"
	"getsturdy.com/api/pkg/users/db"

	"github.com/gin-gonic/gin"
)

func UpdateAvatar(
	logger *zap.Logger,
	repo db.Repository,
	uploader uploader.Uploader,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := auth.UserID(c.Request.Context())
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		form, err := c.MultipartForm()
		if err != nil {
			logger.Warn("failed to parse form", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to update user"})
			return
		}

		file, ok := form.File["file"]
		if !ok {
			logger.Warn("failed to read file", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to update user"})
			return
		}

		if len(file) != 1 {
			logger.Warn("no file found", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to update user"})
			return
		}

		fp, err := file[0].Open()
		if err != nil {
			logger.Warn("failed to open file", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
			return
		}
		defer fp.Close()

		var imgThumb bytes.Buffer
		if err := img.Thumbnail(400, fp, &imgThumb); err != nil {
			logger.Warn("failed to create thumbnail file", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
			return
		}

		key := fmt.Sprintf("%s.png", uuid.New().String())
		avatar, err := uploader.Upload(c.Request.Context(), key, &imgThumb)
		if err != nil {
			logger.Warn("failed to upload avatar", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
			return
		}

		userObj, err := repo.Get(userID)
		if err != nil {
			logger.Warn("failed to get user", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
			return
		}

		userObj.AvatarURL = &avatar.URL
		if err = repo.Update(userObj); err != nil {
			logger.Warn("failed to update user", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
			return
		}

		c.JSON(http.StatusOK, userObj)
	}
}

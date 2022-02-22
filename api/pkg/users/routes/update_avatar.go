package routes

import (
	"bytes"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/img"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"getsturdy.com/api/pkg/users/db"

	"github.com/gin-gonic/gin"
)

func UpdateAvatar(logger *zap.Logger, repo db.Repository, awsSession *session.Session) func(c *gin.Context) {
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
		err = img.Thumbnail(400, fp, &imgThumb)
		if err != nil {
			logger.Warn("failed to create thumbnail file", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		}

		key := fmt.Sprintf("%s.png", uuid.New().String())

		uploader := s3manager.NewUploader(awsSession)
		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String("usercontent.getsturdy.com"),
			Key:    &key,
			Body:   &imgThumb,
		})
		if err != nil {
			logger.Warn("failed to upload avatar to s2", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		}

		url := fmt.Sprintf("https://usercontent.getsturdy.com/%s", key)

		userObj, err := repo.Get(userID)
		if err != nil {
			logger.Warn("failed to get user", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		}

		userObj.AvatarURL = &url
		err = repo.Update(userObj)
		if err != nil {
			logger.Warn("failed to update user", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		}

		c.JSON(http.StatusOK, userObj)
	}
}

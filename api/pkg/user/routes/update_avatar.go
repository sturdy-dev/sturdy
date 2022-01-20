package routes

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/img"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"getsturdy.com/api/pkg/user/db"

	"github.com/gin-gonic/gin"
)

func UpdateAvatar(repo db.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := auth.UserID(c.Request.Context())
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		form, err := c.MultipartForm()
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to update user"})
			return
		}

		file, ok := form.File["file"]
		if !ok {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to update user"})
			return
		}

		if len(file) != 1 {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to update user"})
			return
		}

		sess, err := session.NewSession(
			&aws.Config{
				Region: aws.String("eu-north-1"),
			})
		if err != nil {
			panic(err)
		}

		fp, err := file[0].Open()
		if err != nil {
			panic(err)
		}
		defer fp.Close()

		var imgThumb bytes.Buffer
		err = img.Thumbnail(400, fp, &imgThumb)
		if err != nil {
			panic(err)
		}

		key := fmt.Sprintf("%s.png", uuid.New().String())

		uploader := s3manager.NewUploader(sess)
		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String("usercontent.getsturdy.com"),
			Key:    &key,
			Body:   &imgThumb,
		})

		url := fmt.Sprintf("https://usercontent.getsturdy.com/%s", key)

		userObj, err := repo.Get(userID)
		if err != nil {
			panic(err)
		}
		userObj.AvatarURL = &url
		err = repo.Update(userObj)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, userObj)
	}
}

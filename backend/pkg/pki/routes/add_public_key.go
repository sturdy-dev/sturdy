package routes

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gliderlabs/ssh"

	"mash/pkg/auth"
	"mash/pkg/pki"
	"mash/pkg/pki/db"
)

type AddPublicKeyRequest struct {
	// In AuthorizedKey format
	PublicKey string `json:"public_key" binding:"required"`
}

func AddPublicKey(repo db.Repo) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req AddPublicKeyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		userID, err := auth.UserID(c.Request.Context())
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// If this user already has this key, respond as OK
		if _, err := repo.GetByPublicKeyAndUserID(req.PublicKey, userID); err == nil {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
			return
		}

		upk := pki.UserPublicKey{
			PublicKey: req.PublicKey,
			UserID:    userID,
			AddedAt:   time.Now(),
		}

		err = repo.Create(upk)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not add key"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

type VerifyPublicKeyRequest struct {
	// In what format?
	PublicKey []byte `json:"public_key" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
}

func Verify(repo db.Repo) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req VerifyPublicKeyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		incomingKey, err := ssh.ParsePublicKey(req.PublicKey)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		keys, err := repo.GetKeyByUserID(req.UserID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "not found"})
			return
		}

		for _, key := range keys {
			if key.RevokedAt != nil {
				continue
			}

			parsedAuthorizedKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(key.PublicKey))
			if err != nil {
				log.Println(err)
				continue
			}

			// This user has this key!
			if ssh.KeysEqual(incomingKey, parsedAuthorizedKey) {
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"status": "not found"})
	}
}

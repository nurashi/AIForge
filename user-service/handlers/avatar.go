package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	uploadDir     = "./uploads/avatars"
	maxAvatarSize = 5 << 20
)

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

func UploadAvatar(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt64("user_id")

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxAvatarSize)

		file, header, err := c.Request.FormFile("avatar")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "avatar file is required"})
			return
		}
		defer file.Close()

		ext := strings.ToLower(filepath.Ext(header.Filename))
		if !allowedExtensions[ext] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "allowed formats: jpg, jpeg, png, webp"})
			return
		}

		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			log.Printf("Failed to create upload directory: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save avatar"})
			return
		}

		filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		savePath := filepath.Join(uploadDir, filename)

		dst, err := os.Create(savePath)
		if err != nil {
			log.Printf("Failed to create file %s: %v", savePath, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save avatar"})
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			log.Printf("Failed to write avatar file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save avatar"})
			return
		}

		avatarPath := fmt.Sprintf("/uploads/avatars/%s", filename)

		query := `
			UPDATE user_profiles SET avatar_path = $2
			WHERE user_id = $1
			RETURNING avatar_path`

		var saved string
		err = db.QueryRow(c.Request.Context(), query, userID, avatarPath).Scan(&saved)
		if err != nil {
			log.Printf("Failed to update avatar path for user %d: %v", userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update avatar"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"avatar_path": saved})
	}
}

func DeleteAvatar(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt64("user_id")

		var currentPath string
		err := db.QueryRow(c.Request.Context(),
			`SELECT avatar_path FROM user_profiles WHERE user_id = $1`, userID,
		).Scan(&currentPath)
		if err != nil {
			log.Printf("Failed to get current avatar for user %d: %v", userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete avatar"})
			return
		}

		if currentPath != "" {
			filePath := "." + currentPath
			if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
				log.Printf("Failed to remove avatar file %s: %v", filePath, err)
			}
		}

		_, err = db.Exec(c.Request.Context(),
			`UPDATE user_profiles SET avatar_path = '' WHERE user_id = $1`, userID,
		)
		if err != nil {
			log.Printf("Failed to clear avatar path for user %d: %v", userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete avatar"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "avatar deleted"})
	}
}

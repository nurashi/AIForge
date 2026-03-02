package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nurashi/AIForge/user-service/models"
)

func GetMyProfile(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt64("user_id")

		profile, err := findProfileByUserID(c.Request.Context(), db, userID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				profile, err = createDefaultProfile(c.Request.Context(), db, userID)
				if err != nil {
					log.Printf("Failed to create default profile for user %d: %v", userID, err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create profile"})
					return
				}
			} else {
				log.Printf("Failed to get profile for user %d: %v", userID, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get profile"})
				return
			}
		}

		c.JSON(http.StatusOK, profile)
	}
}

func UpdateMyProfile(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt64("user_id")

		var req models.UpdateProfileRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := findProfileByUserID(c.Request.Context(), db, userID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				_, err = createDefaultProfile(c.Request.Context(), db, userID)
				if err != nil {
					log.Printf("Failed to create default profile for user %d: %v", userID, err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create profile"})
					return
				}
			} else {
				log.Printf("Failed to find profile for user %d: %v", userID, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find profile"})
				return
			}
		}

		query := `
			UPDATE user_profiles
			SET display_name = COALESCE($2, display_name),
				bio = COALESCE($3, bio),
				location = COALESCE($4, location),
				website = COALESCE($5, website)
			WHERE user_id = $1
			RETURNING id, user_id, display_name, bio, location, website, avatar_path, created_at, updated_at`

		var updated models.UserProfile
		err = db.QueryRow(c.Request.Context(), query,
			userID, req.DisplayName, req.Bio, req.Location, req.Website,
		).Scan(
			&updated.ID, &updated.UserID, &updated.DisplayName,
			&updated.Bio, &updated.Location, &updated.Website,
			&updated.AvatarPath, &updated.CreatedAt, &updated.UpdatedAt,
		)
		if err != nil {
			log.Printf("Failed to update profile for user %d: %v", userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
			return
		}

		c.JSON(http.StatusOK, updated)
	}
}

func findProfileByUserID(ctx context.Context, db *pgxpool.Pool, userID int64) (*models.UserProfile, error) {
	query := `
		SELECT id, user_id, display_name, bio, location, website, avatar_path, created_at, updated_at
		FROM user_profiles
		WHERE user_id = $1`

	var p models.UserProfile
	err := db.QueryRow(ctx, query, userID).Scan(
		&p.ID, &p.UserID, &p.DisplayName, &p.Bio,
		&p.Location, &p.Website, &p.AvatarPath,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func createDefaultProfile(ctx context.Context, db *pgxpool.Pool, userID int64) (*models.UserProfile, error) {
	query := `
		INSERT INTO user_profiles (user_id)
		VALUES ($1)
		RETURNING id, user_id, display_name, bio, location, website, avatar_path, created_at, updated_at`

	var p models.UserProfile
	err := db.QueryRow(ctx, query, userID).Scan(
		&p.ID, &p.UserID, &p.DisplayName, &p.Bio,
		&p.Location, &p.Website, &p.AvatarPath,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

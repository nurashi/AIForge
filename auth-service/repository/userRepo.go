package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5" 
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nurashi/AIForge/auth-service/models" 
)

type UserRepository interface {
	FindUserByGoogleID(ctx context.Context, googleID string) (*models.User, error)
	CreateUser(ctx context.Context, userInfo *models.GoogleUserInfo) (*models.User, error)
}

type Userrepository struct {
	db *pgxpool.Pool 
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &Userrepository{db: db}
}

func (r *Userrepository) FindUserByGoogleID(ctx context.Context, googleID string) (*models.User, error) {
	query := `
        SELECT id, google_id, email, name, avatar_url, created_at, updated_at
        FROM users
        WHERE google_id = $1` 

	user := &models.User{} 

	err := r.db.QueryRow(ctx, query, googleID).Scan(
		&user.ID,
		&user.GoogleID,
		&user.Email,
		&user.Name,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		log.Printf("Repository ERROR: finding user by Google ID %s: %v", googleID, err)
		return nil, fmt.Errorf("database query error: %w", err)
	}

	return user, nil
}

func (r *Userrepository) CreateUser(ctx context.Context, userInfo *models.GoogleUserInfo) (*models.User, error) {
	query := `
        INSERT INTO users (google_id, email, name, avatar_url)
        VALUES ($1, $2, $3, $4)
        RETURNING id, google_id, email, name, avatar_url, created_at, updated_at`

	newUser := &models.User{} 

	err := r.db.QueryRow(ctx, query, userInfo.Sub,       
		userInfo.Email,   
		userInfo.Name,     
		userInfo.Picture,   
	).Scan(
		&newUser.ID,
		&newUser.GoogleID,
		&newUser.Email,
		&newUser.Name,
		&newUser.AvatarURL,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	)

	if err != nil {
		log.Printf("Repository ERROR: Failed creating user with email %s: %v", userInfo.Email, err)
		return nil, fmt.Errorf("database insert error: %w", err)
	}

	log.Printf("Repository: Successfully createded user ID %d for email %s", newUser.ID, newUser.Email)
	return newUser, nil
}

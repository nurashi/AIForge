package models

import "time"

type UserProfile struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	DisplayName string    `json:"display_name"`
	Bio         string    `json:"bio"`
	Location    string    `json:"location"`
	Website     string    `json:"website"`
	AvatarPath  string    `json:"avatar_path"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateProfileRequest struct {
	DisplayName *string `json:"display_name" binding:"omitempty,max=100"`
	Bio         *string `json:"bio" binding:"omitempty,max=2000"`
	Location    *string `json:"location" binding:"omitempty,max=255"`
	Website     *string `json:"website" binding:"omitempty,max=255"`
}

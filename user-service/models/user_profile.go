package models

import (
	"time"
)

type UserProfile struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`    
	Bio       string    `json:"bio"`        
	Location  string    `json:"location"`  
	LinkedIn  string    `json:"linkedin"`  
	GitHub    string    `json:"github"`     
	CreatedAt time.Time `json:"created_at"` 
	UpdatedAt time.Time `json:"updated_at"` 

}

type PersonalProject struct {
	ID           int64     `json:"id"`
	UserProfileID int64    `json:"user_profile_id"` // fk for UserProfile
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Technologies []string `json:"technologies"` // JSONB in DB
	Link         string    `json:"link"`
	StartDate    time.Time `json:"start_date"`
	EndDate      *time.Time `json:"end_date"` 
	IsFeatured   bool      `json:"is_featured"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}


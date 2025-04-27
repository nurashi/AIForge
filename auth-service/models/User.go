package models

import (
	"time"
)

type User struct {
	ID        int64     `json:"id"`         
	GoogleID  string    `json:"-"`       
	Email     string    `json:"email"`   
	Name      string    `json:"name"`       
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"` 
	UpdatedAt time.Time `json:"updated_at"`
}

type GoogleUserInfo struct {
	Sub           string `json:"sub"`         
	Name          string `json:"name"`          
	GivenName     string `json:"given_name"`   
	FamilyName    string `json:"family_name"`   
	Profile       string `json:"profile"`       
	Picture       string `json:"picture"`      
	Email         string `json:"email"`          
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`        
	Hd            string `json:"hd"`            
}


package repository

import (
	"time"
)

// SocialRequest represents the request body for creating a new social media entry.
type SocialRequest struct {
	Name           string `json:"name" example:"Facebook"`
	SocialMediaUrl string `json:"social_media_url" example:"https://www.facebook.com/myusername"`
}

// SocialCreateResponse represents the response body for creating a new social media entry.
type SocialCreateResponse struct {
	Id             uint       `json:"id" example:"1"`
	Name           string     `json:"name" example:"Facebook"`
	SocialMediaUrl string     `json:"social_media_url" example:"https://www.facebook.com/myusername"`
	UserId         uint       `json:"user_id,omitempty" example:"1"`
	CreatedAt      *time.Time `json:"created_at,omitempty" example:"2023-04-15T14:30:00Z"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty" example:"2023-04-15T14:30:00Z"`
}

// SocialGetResponse represents the response body for getting a list of social media entries.
type SocialGetResponse struct {
	Socials []SocialData `json:"social_medias"`
}

// SocialData represents a single social media entry in the response body.
type SocialData struct {
	Id             uint               `json:"id" example:"1"`
	Name           string             `json:"name" example:"Facebook"`
	SocialMediaUrl string             `json:"social_media_url" example:"https://www.facebook.com/myusername"`
	CreatedAt      *time.Time         `json:"created_at" example:"2023-04-15T14:30:00Z"`
	UpdatedAt      *time.Time         `json:"updated_at" example:"2023-04-15T14:30:00Z"`
	User           UserSocialResponse `json:"user"`
}

// UserSocialResponse represents the user data associated with a social media entry in the response body.
type UserSocialResponse struct {
	Id       uint   `json:"id" example:"1"`
	Username string `json:"username" example:"john_doe"`
}

package repository

import "time"

// CommentRequest represents the request body for creating a comment
type CommentRequest struct {
	Message string `json:"message"`
	PhotoId uint   `json:"photo_id"`
}

// CommentCreateResponse represents the response body for creating a comment
type CommentCreateResponse struct {
	Id        uint       `json:"id"`
	Message   string     `json:"message"`
	PhotoId   uint       `json:"photo_id"`
	UserId    uint       `json:"user_id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// CommentGetResponse represents the response body for getting multiple comments
type CommentGetResponse struct {
	Comments []CommentData `json:"photos"`
}

// CommentData represents the data of a comment
type CommentData struct {
	Id        uint                 `json:"id"`
	Message   string               `json:"message"`
	PhotoId   uint                 `json:"photo_id"`
	User      UserCommentResponse  `json:"user"`
	Photos    PhotoCommentResponse `json:"photos"`
	CreatedAt *time.Time           `json:"created_at"`
	UpdatedAt *time.Time           `json:"updated_at"`
}

// UserCommentResponse represents the response body for a user associated with a comment
type UserCommentResponse struct {
	Id       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// PhotoCommentResponse represents the response body for a photo associated with a comment
type PhotoCommentResponse struct {
	Id       uint   `json:"id"`
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url"`
	UserId   uint   `json:"user_id"`
}

// CommentUpdateRequest represents the request body for updating a comment
type CommentUpdateRequest struct {
	Message string `json:"message"`
}

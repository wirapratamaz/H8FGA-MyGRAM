package repository

import "time"

// swagger:parameters createPhotoRequest
type PhotoRequest struct {
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url"`
}

// swagger:response photoCreateResponse
type PhotoCreateResponse struct {
	Id        uint       `json:"id"`
	Title     string     `json:"title"`
	Caption   string     `json:"caption"`
	PhotoUrl  string     `json:"photo_url"`
	UserId    uint       `json:"user_id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// swagger:response photoGetDataResponse
type PhotoGetResponse struct {
	Photos []PhotoData `json:"photos"`
}

// Data foto
// swagger:model photoData
type PhotoData struct {
	Id        uint              `json:"id"`
	Title     string            `json:"title"`
	Caption   string            `json:"caption"`
	PhotoUrl  string            `json:"photo_url"`
	User      UserPhotoResponse `json:"user"`
	CreatedAt *time.Time        `json:"created_at"`
	UpdatedAt *time.Time        `json:"updated_at"`
}

// Data user pada foto
// swagger:model userPhotoResponse
type UserPhotoResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

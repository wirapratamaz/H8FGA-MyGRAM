package repository

import "time"

// Objek Response saat user berhasil dibuat
// swagger:response userCreateResponse
type UserCreateResponse struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
}

// Objek Response saat user berhasil login
// swagger:response userLoginResponse
type UserLoginResponse struct {
	Token string `json:"token"`
}

// Objek Request saat meng-update informasi user
// swagger:parameters userUpdateRequest
type UserUpdateRequest struct {
	Email    string `json:"email" valid:"email~Invalid format email"`
	Username string `json:"username"`
}

// Objek Response saat informasi user berhasil di-update
// swagger:response userUpdateResponse
type UserUpdateResponse struct {
	Id        uint       `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Age       int        `json:"age"`
	UpdatedAt *time.Time `json:"updated_at"`
}

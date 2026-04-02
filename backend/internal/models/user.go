package models

import (
	"time"
)

type User struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	PasswordHash string     `json:"-"`
	FullName     string     `json:"full_name"`
	Role         string     `json:"role"`
	NIP          *string    `json:"nip,omitempty"`
	Email        *string    `json:"email,omitempty"`
	Phone        *string    `json:"phone,omitempty"`
	IsActive     bool       `json:"is_active"`
	LastLogin    *time.Time `json:"last_login,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type CreateUserRequest struct {
	Username string  `json:"username" binding:"required"`
	Password string  `json:"password" binding:"required,min=6"`
	FullName string  `json:"full_name" binding:"required"`
	Role     string  `json:"role" binding:"required,oneof=ADMIN TEACHER"`
	NIP      *string `json:"nip"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
}

type UpdateUserRequest struct {
	Username *string `json:"username"`
	FullName *string `json:"full_name"`
	Role     *string `json:"role" binding:"omitempty,oneof=ADMIN TEACHER"`
	NIP      *string `json:"nip"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	IsActive *bool   `json:"is_active"`
}

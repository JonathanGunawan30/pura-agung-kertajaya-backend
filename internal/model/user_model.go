package model

import "time"

type UserResponse struct {
	ID        int       `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type RegisterUserRequest struct {
	Name     string `json:"name" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type LoginUserRequest struct {
	Email          string `json:"email" validate:"required,email,max=100"`
	Password       string `json:"password" validate:"required,max=100"`
	RecaptchaToken string `json:"recaptcha_token,omitempty"`
}

type UpdateUserRequest struct {
	ID       int    `json:"-"`
	Name     string `json:"name,omitempty" validate:"max=100"`
	Password string `json:"password,omitempty" validate:"max=100"`
}

type GetUserRequest struct {
	ID int `json:"id" validate:"required"`
}

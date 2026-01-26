package account

import (
	"time"
)

type Account struct {
	ID       string `json:"id" validate:"required,uuidv4"`
	Name     string `json:"name" validate:"min=3,max=50"`
	Email    string `json:"email" validate:"required,email,normalizeemail"`
	Password string `json:"-" validate:"required,min=8,max=50"`
}

type AuthenticatedResponse struct {
	Account      *Account `json:"account"`
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
}

type Session struct {
	ID           string    `json:"id"`
	AccountID    string    `json:"account_id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	IsRevoked    bool      `json:"is_revoked"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

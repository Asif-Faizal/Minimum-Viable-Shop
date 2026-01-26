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
	DeviceID     string    `json:"device_id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	IsRevoked    bool      `json:"is_revoked"`
}

type DeviceInfo struct {
	ID              string    `json:"id"`
	SessionID       string    `json:"session_id"`
	DeviceType      string    `json:"device_type"`
	DeviceModel     string    `json:"device_model"`
	DeviceOS        string    `json:"device_os"`
	DeviceOSVersion string    `json:"device_os_version"`
	IPAddress       string    `json:"ip_address"`
	UserAgent       string    `json:"user_agent"`
	CreatedAt       time.Time `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

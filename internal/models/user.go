package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Phone     string     `json:"phone,omitempty"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Role      string     `json:"role,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Session struct {
	RefreshToken string    `json:"refreshToken" bson:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt" bson:"expiresAt"`
}

type SignInRequest struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdatePasswordRequest struct {
	UserID      int    `json:"user_id"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type SignUpResponse struct {
	User             User   `json:"user"`
	VerificationCode string `json:"verification_code,omitempty"`
}
type SignUpResponse1 struct {
	User             User   `json:"user"`
	VerificationCode string `json:"verification_code"`
}

type DuplicateCheckRequest struct {
	Phone string `json:"phone"`
	Email string `json:"email"`
}

type VerificationCodeEntry struct {
	ID        int       `json:"id"`
	Phone     string    `json:"phone"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
}

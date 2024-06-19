package entity

import "time"

type User struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Email           string    `json:"email" gorm:"unique;not null"`
	Password        string    `json:"password"`
	OTP             string    `json:"otp"`
	OTPExpiresAt    time.Time `json:"otp_expires_at"`
	EmailVerifiedAt time.Time `json:"email_verified_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

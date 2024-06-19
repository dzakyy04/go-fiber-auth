package entity

import "time"

type User struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	Name            string     `json:"name"`
	Email           string     `json:"email" gorm:"unique;not null"`
	Password        string     `json:"password"`
	OTP             *string    `json:"otp" gorm:"default:null"`
	OTPExpiresAt    *time.Time `json:"otp_expires_at" gorm:"default:null"`
	EmailVerifiedAt *time.Time `json:"email_verified_at" gorm:"default:null"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

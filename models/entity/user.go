package entity

import "time"

type User struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	Name            string     `json:"name"`
	Email           string     `json:"email" gorm:"unique;not null"`
	Password        string     `json:"-" gorm:"column:password;not null"`
	OTP             *string    `json:"-" gorm:"column:otp;default:null"`
	OTPExpiresAt    *time.Time `json:"-" gorm:"column:otp_expires_at;default:null"`
	EmailVerifiedAt *time.Time `json:"email_verified_at" gorm:"default:null"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

package request

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type EmailVerificationRequest struct {
	Email string `json:"email"`
}

type VerifyEmailRequest struct {
	Email string `json:"email"`
	OTP	string `json:"otp"`
}

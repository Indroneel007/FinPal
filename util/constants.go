package util

import (
	//"os"
	"time"
)

const (
	authTokenExp       = time.Minute * 15
	refreshTokenExp    = time.Hour * 24 * 30
	blacklistKeyPrefix = "blacklisted:"
	OtpKeyPrefix       = "password-reset:"
	otpExp             = time.Minute * 10
	otpCharSet         = "1234567890"
	emailTemplate      = "To: %s\r\n" +
		"Subject: FinPal Password Reset\r\n" +
		"\r\n" +
		"Your OTP for password reset is %s\r\n"

	// public because needed for testing
	OTPLength = 4
)

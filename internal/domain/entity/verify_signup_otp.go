package entity

import "time"

type OTP struct {
	ID        string
	Phone     string
	Code      string
	ExpiresAt time.Time
	Verified  bool
}

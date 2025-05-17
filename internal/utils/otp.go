package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateOTP generates a random 6-digit OTP code as string
func GenerateOTP() string {
	// Create a new source with the current time
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	// Generate a random number between 100000 and 999999
	return fmt.Sprintf("%06d", r.Intn(900000)+100000)
}

// GenerateOTPExpiry returns a time 15 minutes from now
func GenerateOTPExpiry() time.Time {
	return time.Now().Add(15 * time.Minute)
}

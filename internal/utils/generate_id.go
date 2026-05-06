package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateSixDigitID generates a random six-digit ID.
func GenerateSixDigitID() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

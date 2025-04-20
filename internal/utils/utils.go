package utils

import (
	"time"
)

// GenerateRefNum creates a unique reference number with a timestamp and random component
func GenerateRefNum() string {
	timestamp := time.Now().Format("20060102150405")
	return timestamp
}

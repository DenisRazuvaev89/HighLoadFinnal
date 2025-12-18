package utils

import (
	"log"
	"time"
)

func LogUserAction(action string, userID int) {
	timestamp := time.Now().Format(time.RFC3339)
	log.Printf("[AUDIT] %s - Action: %s - UserID: %d", timestamp, action, userID)
}

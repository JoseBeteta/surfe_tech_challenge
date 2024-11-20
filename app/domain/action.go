package domain

import "time"

type Action struct {
	ID         int       `json:"id"`
	Type       string    `json:"type"`
	UserID     int       `json:"userId"`
	TargetUser int       `json:"targetUser,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
}

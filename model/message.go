package model

import "time"

// Message data
type Message struct {
	ID        int       `json:"id"`
	User      User      `json:"user"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	Room      string    `json:"room"`
}

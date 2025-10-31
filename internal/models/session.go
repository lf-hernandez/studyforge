package models

import "time"

// Session represents a user session
type Session struct {
	ID           string    `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	LastAccessed time.Time `json:"last_accessed"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	IsActive     bool      `json:"is_active"`
}

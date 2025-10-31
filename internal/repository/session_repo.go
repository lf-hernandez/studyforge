package repository

import (
	"database/sql"
	"fmt"
	"time"

	"studyforge/internal/models"
)

// SessionRepository handles session database operations
type SessionRepository struct {
	db *sql.DB
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Create creates a new session
func (r *SessionRepository) Create(session *models.Session) error {
	query := `
		INSERT INTO sessions (id, created_at, last_accessed, ip_address, user_agent, is_active)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query,
		session.ID,
		session.CreatedAt,
		session.LastAccessed,
		session.IPAddress,
		session.UserAgent,
		session.IsActive,
	)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

// GetByID retrieves a session by ID
func (r *SessionRepository) GetByID(id string) (*models.Session, error) {
	query := `
		SELECT id, created_at, last_accessed, ip_address, user_agent, is_active
		FROM sessions
		WHERE id = ? AND is_active = TRUE
	`
	session := &models.Session{}
	err := r.db.QueryRow(query, id).Scan(
		&session.ID,
		&session.CreatedAt,
		&session.LastAccessed,
		&session.IPAddress,
		&session.UserAgent,
		&session.IsActive,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	return session, nil
}

// UpdateLastAccessed updates the last accessed timestamp
func (r *SessionRepository) UpdateLastAccessed(id string) error {
	query := `UPDATE sessions SET last_accessed = ? WHERE id = ?`
	_, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update last accessed: %w", err)
	}
	return nil
}

// Deactivate deactivates a session
func (r *SessionRepository) Deactivate(id string) error {
	query := `UPDATE sessions SET is_active = FALSE WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to deactivate session: %w", err)
	}
	return nil
}

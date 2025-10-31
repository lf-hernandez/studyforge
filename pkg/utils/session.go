package utils

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"studyforge/internal/models"
	"studyforge/internal/repository"

	"github.com/google/uuid"
)

type contextKey string

const SessionContextKey contextKey = "session"

// SessionManager handles session creation and validation
type SessionManager struct {
	sessionRepo *repository.SessionRepository
}

// NewSessionManager creates a new session manager
func NewSessionManager(sessionRepo *repository.SessionRepository) *SessionManager {
	return &SessionManager{
		sessionRepo: sessionRepo,
	}
}

// Middleware provides session management for HTTP handlers
func (sm *SessionManager) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to get session from cookie
		cookie, err := r.Cookie("session_id")
		var session *models.Session

		if err == nil && cookie.Value != "" {
			// Try to load existing session
			session, err = sm.sessionRepo.GetByID(cookie.Value)
			if err == nil {
				// Update last accessed time
				if err := sm.sessionRepo.UpdateLastAccessed(session.ID); err != nil {
					log.Printf("Failed to update session last accessed: %v", err)
				}
			}
		}

		// Create new session if needed
		if session == nil {
			session, err = sm.createSession(r)
			if err != nil {
				log.Printf("Failed to create session: %v", err)
				http.Error(w, "Session error", http.StatusInternalServerError)
				return
			}

			// Set cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    session.ID,
				Path:     "/",
				MaxAge:   86400, // 24 hours
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			})
		}

		// Add session to context
		ctx := context.WithValue(r.Context(), SessionContextKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// createSession creates a new session
func (sm *SessionManager) createSession(r *http.Request) (*models.Session, error) {
	session := &models.Session{
		ID:           uuid.New().String(),
		CreatedAt:    time.Now(),
		LastAccessed: time.Now(),
		IPAddress:    r.RemoteAddr,
		UserAgent:    r.UserAgent(),
		IsActive:     true,
	}

	if err := sm.sessionRepo.Create(session); err != nil {
		return nil, err
	}

	log.Printf("Created new session: %s", session.ID)
	return session, nil
}

// GetSessionFromContext retrieves session from request context
func GetSessionFromContext(ctx context.Context) (*models.Session, error) {
	session, ok := ctx.Value(SessionContextKey).(*models.Session)
	if !ok {
		return nil, fmt.Errorf("no session in context")
	}
	return session, nil
}

package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

// SessionData holds the user ID and expiration for a session
type SessionData struct {
	UserID    string // UserID can be string or int64, depending on your user model
	ExpiresAt time.Time
}

// Store is a simple in-memory session store
type Store struct {
	mu       sync.RWMutex
	sessions map[string]SessionData // token -> SessionData
}

// NewStore creates a new in-memory session store
func NewStore() *Store {
	s := &Store{
		sessions: make(map[string]SessionData),
	}
	// Periodically clean up expired sessions (optional, but good practice)
	go s.cleanupExpiredSessions()
	return s
}

// GenerateToken creates a new random session token
func GenerateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// CreateSession stores a new session
func (s *Store) CreateSession(userID string, duration time.Duration) (string, error) {
	token, err := GenerateToken()
	if err != nil {
		return "", err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[token] = SessionData{
		UserID:    userID,
		ExpiresAt: time.Now().Add(duration),
	}
	return token, nil
}

// ValidateToken checks if a token is valid and returns the UserID
// You would pass this function or the Store itself to your AuthMiddleware
func (s *Store) ValidateToken(tokenString string) (userID string, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sessionData, exists := s.sessions[tokenString]
	if !exists {
		return "", errors.New("invalid or non-existent session token")
	}

	if time.Now().After(sessionData.ExpiresAt) {
		// Token expired, remove it (lazy cleanup)
		s.mu.RUnlock() // Release read lock
		s.mu.Lock()    // Acquire write lock
		delete(s.sessions, tokenString)
		s.mu.Unlock() // Release write lock
		s.mu.RLock()  // Re-acquire read lock (though not strictly needed before return)
		return "", errors.New("session token expired")
	}

	return sessionData.UserID, nil
}

// DeleteSession removes a session (for logout)
func (s *Store) DeleteSession(tokenString string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, tokenString)
}

func (s *Store) cleanupExpiredSessions() {
	for {
		time.Sleep(10 * time.Minute) // Check every 10 minutes
		s.mu.Lock()
		for token, data := range s.sessions {
			if time.Now().After(data.ExpiresAt) {
				delete(s.sessions, token)
			}
		}
		s.mu.Unlock()
	}
}

// Define an interface for cleaner dependency injection if preferred
type SessionValidator interface {
	ValidateToken(tokenString string) (userID string, err error)
}

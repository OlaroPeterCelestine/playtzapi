package auth

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"
)

// Session represents a user session
type Session struct {
	ID        string
	UserID    string
	Username  string
	Email     string
	RoleID    string
	RoleName  string
	LastActivity time.Time
	ExpiresAt    time.Time
}

// SessionStore manages user sessions
type SessionStore struct {
	sessions map[string]*Session
	mutex    sync.RWMutex
	timeout  time.Duration
}

var store *SessionStore
var once sync.Once

// GetSessionStore returns the singleton session store
func GetSessionStore() *SessionStore {
	once.Do(func() {
		store = &SessionStore{
			sessions: make(map[string]*Session),
			timeout:  10 * time.Minute, // 10 minute inactivity timeout
		}
		// Start cleanup goroutine
		go store.cleanup()
	})
	return store
}

// CreateSession creates a new session for a user
func (s *SessionStore) CreateSession(userID, username, email, roleID, roleName string) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	sessionID := generateSessionID()
	now := time.Now()
	
	session := &Session{
		ID:          sessionID,
		UserID:      userID,
		Username:    username,
		Email:       email,
		RoleID:      roleID,
		RoleName:    roleName,
		LastActivity: now,
		ExpiresAt:    now.Add(s.timeout),
	}

	s.sessions[sessionID] = session
	return sessionID
}

// GetSession retrieves a session by ID and updates last activity
func (s *SessionStore) GetSession(sessionID string) (*Session, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return nil, false
	}

	// Check if session expired
	if time.Now().After(session.ExpiresAt) {
		delete(s.sessions, sessionID)
		return nil, false
	}

	// Update last activity
	session.LastActivity = time.Now()
	session.ExpiresAt = session.LastActivity.Add(s.timeout)
	
	return session, true
}

// DeleteSession removes a session
func (s *SessionStore) DeleteSession(sessionID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.sessions, sessionID)
}

// cleanup removes expired sessions periodically
func (s *SessionStore) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mutex.Lock()
		now := time.Now()
		for id, session := range s.sessions {
			// Remove if expired or inactive for more than timeout
			if now.After(session.ExpiresAt) || now.Sub(session.LastActivity) > s.timeout {
				delete(s.sessions, id)
			}
		}
		s.mutex.Unlock()
	}
}

// generateSessionID generates a random session ID
func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}


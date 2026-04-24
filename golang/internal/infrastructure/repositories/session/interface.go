package sessionrepository

import "time"

// SessionRepository defines the interface for SessionRepository use case
type SessionRepository interface {
	SetSession(key string, value interface{}, expiration time.Duration) error
	GetSession(key string) (string, error)
	DeleteSession(key string) error

	GetAuthCodePrefix() string
	SaveAuthCode(code, sessionID string) error
	GetAuthCode(code string) (string, error)
	DeleteAuthCode(code string) error
}

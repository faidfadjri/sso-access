package sessionrepository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)


const AuthCodePrefix = "authcode:"

type sessionRepository struct {
	redisClient *redis.Client
}

func NewSessionRepository(redisClient *redis.Client) SessionRepository {
	return &sessionRepository{
		redisClient: redisClient,
	}
}

func (r *sessionRepository) SetSession(key string, value interface{}, expiration time.Duration) error {
	return r.redisClient.Set(context.Background(), key, value, expiration).Err()
}

func (r *sessionRepository) GetSession(key string) (string, error) {
	return r.redisClient.Get(context.Background(), key).Result()
}

func (r *sessionRepository) DeleteSession(key string) error {
	return r.redisClient.Del(context.Background(), key).Err()
}

// ----- AUTHORIZATION CODE ----- //

func (r *sessionRepository) GetAuthCodePrefix() string {
	return AuthCodePrefix
}

func (r *sessionRepository) SaveAuthCode(code, sessionID string) error {
	key := AuthCodePrefix + code
	return r.SetSession(key, sessionID, 5*time.Minute)
}

func (r *sessionRepository) GetAuthCode(code string) (string, error) {
	key := AuthCodePrefix + code
	return r.GetSession(key)
}

func (r *sessionRepository) DeleteAuthCode(code string) error {
	key := AuthCodePrefix + code
	return r.DeleteSession(key)
}
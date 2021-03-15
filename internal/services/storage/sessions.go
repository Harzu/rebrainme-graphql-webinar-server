package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"webinar/graphql/server/internal/entities"
	"webinar/graphql/server/internal/system/database/redis"
)

const sessionTTL = 30 * time.Minute

type Session interface {
	Set(ctx context.Context, sessionKey string, session entities.Session) error
	Get(ctx context.Context, sessionKey string) (entities.Session, error)
}

type sessionStorage struct {
	redisClient redis.Repository
}

func NewSessionStorage(redisClient redis.Repository) Session {
	return &sessionStorage{
		redisClient: redisClient,
	}
}

func (s *sessionStorage) Set(ctx context.Context, sessionKey string, session entities.Session) error {
	rawData, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	return s.redisClient.GetConnection().Set(ctx, sessionKey, rawData, sessionTTL).Err()
}

func (s *sessionStorage) Get(ctx context.Context, sessionKey string) (result entities.Session, err error) {
	data, err := s.redisClient.GetConnection().Get(ctx, sessionKey).Result()
	if err != nil {
		return result, fmt.Errorf("failed to get session from redis: %w", err)
	}

	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return result, fmt.Errorf("failed to parse session: %w", err)
	}

	return
}

package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type sessionRepository struct {
	client *redis.Client
}

func NewSessionRepository(client *redis.Client) SessionRepository {
	return &sessionRepository{client: client}
}

func key(sessionID string) string {
	return "session:" + sessionID
}

func (s *sessionRepository) Create(ctx context.Context, sessionID string, userID int64, ttl time.Duration) error {
	if err := s.client.Set(ctx, key(sessionID), strconv.FormatInt(userID, 10), ttl).Err(); err != nil {
		return fmt.Errorf("set session: %w", err)
	}

	return nil
}

func (s *sessionRepository) Exists(ctx context.Context, sessionID string) (bool, error) {
	count, err := s.client.Exists(ctx, key(sessionID)).Result()
	if err != nil {
		return false, fmt.Errorf("check session: %w", err)
	}

	return count > 0, nil
}

func (s *sessionRepository) Delete(ctx context.Context, sessionID string) error {
	if err := s.client.Del(ctx, key(sessionID)).Err(); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}

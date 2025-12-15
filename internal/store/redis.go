package store

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	cli *redis.Client
}

func (r *RedisStore) Enqueue(ctx context.Context, jobID string) error {
	if r.cli == nil {
		return errors.New("redis not initialized")
	}
	return r.cli.LPush(ctx, "job_queue", jobID).Err()
}

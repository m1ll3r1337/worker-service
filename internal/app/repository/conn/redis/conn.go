package rcredis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/m1ll3r1337/worker-service/internal/app/config/section"
)

// Client встраивает *redis.Client — наружу доступны все команды go-redis напрямую.
type Client struct {
	*redis.Client
}

func NewClient(ctx context.Context, cfg section.RepositoryRedis) (*Client, error) {
	log.Info().Str("address", cfg.Address).Int("db", cfg.DB).Msg("Connecting to Redis")

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	if err := redisotel.InstrumentTracing(rdb); err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("instrument redis tracing: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := rdb.Ping(pingCtx).Err(); err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	log.Info().Str("address", cfg.Address).Msg("Redis connected")

	return &Client{
		rdb,
	}, nil
}

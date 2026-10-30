package rcurrency

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/m1ll3r1337/worker-service/internal/app/entity"
	"github.com/m1ll3r1337/worker-service/internal/app/repository"
)

const rateKeyPrefix = "rates:"

// repoRedis — кэш курсов валют поверх Redis.
type repoRedis struct {
	client   *redis.Client
	cacheTTL time.Duration
}

func NewRepoFromRedis(client *redis.Client, cacheTTL time.Duration) repository.CurrencyRate {
	return &repoRedis{client: client, cacheTTL: cacheTTL}
}

func (r *repoRedis) buildKey(from, to string) string {
	return fmt.Sprintf("%s%s:%s", rateKeyPrefix, from, to)
}

func (r *repoRedis) GetRate(ctx context.Context, from, to string) (float64, error) {
	key := r.buildKey(from, to)
	res, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, entity.ErrRateNotFound
		}

		return 0, err
	}

	f, err := strconv.ParseFloat(res, 64)
	if err != nil {
		return 0, fmt.Errorf("parse cached rate %q: %w", res, err)
	}

	return f, nil
}

func (r *repoRedis) SetRate(ctx context.Context, from, to string, rate float64) error {
	s := strconv.FormatFloat(rate, 'f', -1, 64)
	key := r.buildKey(from, to)
	err := r.client.Set(ctx, key, s, r.cacheTTL).Err()
	return err
}

func (r *repoRedis) SetRates(ctx context.Context, from string, rates map[string]float64) error {
	pipe := r.client.Pipeline()
	for k, v := range rates {
		key := r.buildKey(from, k)
		pipe.Set(ctx, key, v, r.cacheTTL)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("set rates cache: %w", err)
	}

	log.Info().Ctx(ctx).Str("from", from).Int("count", len(rates)).Msg("Курсы валют сохранены в кэш")

	return nil
}

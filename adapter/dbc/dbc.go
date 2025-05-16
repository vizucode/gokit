package dbc

import (
	"context"
	"database/sql"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// sqlDBc is instance for database/sql connection
type SqlDBc struct {
	DB *sql.DB
}

// GormDBc is instance for gorm connection
type GormDBc struct {
	DB *gorm.DB
}

// RedisDBc is instance for redis connection
type RedisDBc struct {
	DB CacheClient
}

// Redis standard interface for abstract redis
type CacheClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Keys(ctx context.Context, pattern string) *redis.StringSliceCmd
}

package dbc

import (
	"database/sql"

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
	DB      *redis.Client
	Cluster *redis.ClusterClient
}

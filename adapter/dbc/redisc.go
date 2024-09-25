package dbc

import (
	"context"
	"crypto/tls"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/vizucode/gokit/utils/env"
)

type OptionRedis func(o *optionRedis)

type optionRedis struct {
	uri         string
	serviceName string
	// secureTLS defines the TLS configuration for secure communication.
	// It is highly recommended to use TLS for enhanced security.
	secureTLS *tls.Config
	// set minimal connection on pool
	minPoolConnection uint
	// set maximum client can connect
	maxPoolConnection uint
	// set maximum idle connection on pool
	maxIdleConnectionDuration time.Duration
	// set wait connection  time out when pool is all used
	waitPoolConnectionDuration time.Duration
	// set maxLifeTime connection will close relative after connection used
	maxLifeTimeConnection time.Duration
}

func defaultRedisConnection() optionRedis {
	return optionRedis{
		uri:                        env.GetString("DB_REDIS_URI"),
		serviceName:                "default-service-name",
		minPoolConnection:          1,
		maxPoolConnection:          100,
		maxIdleConnectionDuration:  time.Minute * 1,
		waitPoolConnectionDuration: time.Minute * 1,
		secureTLS:                  nil,
	}
}

func NewRedisConnection(options ...OptionRedis) *RedisDBc {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn := defaultRedisConnection()
	for _, option := range options {
		option(&conn)
	}

	opt, err := goredis.ParseURL(conn.uri)
	if err != nil {
		panic(err)
	}

	opt.MaxActiveConns = int(conn.maxPoolConnection)
	opt.MinIdleConns = int(conn.minPoolConnection)
	opt.PoolTimeout = conn.waitPoolConnectionDuration
	opt.ConnMaxIdleTime = conn.maxIdleConnectionDuration
	opt.ConnMaxLifetime = conn.maxLifeTimeConnection
	opt.TLSConfig = conn.secureTLS

	rediscli := goredis.NewClient(opt)

	cmd := rediscli.Ping(ctx)
	if cmd.Err() != nil {
		panic(cmd.Err())
	}

	return &RedisDBc{
		DB: rediscli,
	}
}

func SetRedisURIConnection(uri string) OptionRedis {
	return func(o *optionRedis) {
		o.uri = uri
	}
}

func SetRedisServiceName(serviceName string) OptionRedis {
	return func(o *optionRedis) {
		o.serviceName = serviceName
	}
}

func SetRedisMinPoolConnection(minPoolConnection uint) OptionRedis {
	return func(o *optionRedis) {
		o.minPoolConnection = minPoolConnection
	}
}

func SetRedisMaxPoolConnection(maxPoolConnection uint) OptionRedis {
	return func(o *optionRedis) {
		o.maxPoolConnection = maxPoolConnection
	}
}

func SetRedisMaxIdleConnectionDuration(maxIdleConnectionDuration time.Duration) OptionRedis {
	return func(o *optionRedis) {
		o.maxIdleConnectionDuration = maxIdleConnectionDuration
	}
}

func SetRedisWaitPoolConnectionDuration(waitPoolConnectionDuration time.Duration) OptionRedis {
	return func(o *optionRedis) {
		o.waitPoolConnectionDuration = waitPoolConnectionDuration
	}
}

func SetRedisMaxLifeTimeConnection(maxLifeTimeConnection time.Duration) OptionRedis {
	return func(o *optionRedis) {
		o.maxLifeTimeConnection = maxLifeTimeConnection
	}
}

func SetRedisSecureTLS(secureTLS *tls.Config) OptionRedis {
	return func(o *optionRedis) {
		o.secureTLS = secureTLS
	}
}

package dbc

import (
	"context"
	"crypto/tls"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/vizucode/gokit/utils/env"
)

type OptionRedisCluster func(o *optionRedisCluster)

type optionRedisCluster struct {
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

func defaultRedisClusterConnection() optionRedisCluster {
	return optionRedisCluster{
		uri:                        env.GetString("DB_REDIS_URI"),
		serviceName:                "default-service-name",
		minPoolConnection:          1,
		maxPoolConnection:          100,
		maxIdleConnectionDuration:  time.Minute * 1,
		waitPoolConnectionDuration: time.Minute * 1,
		secureTLS:                  nil,
	}
}

func NewRedisClusterConnection(options ...OptionRedisCluster) *RedisDBc {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn := defaultRedisClusterConnection()
	for _, option := range options {
		option(&conn)
	}

	opt, err := goredis.ParseClusterURL(conn.uri)
	if err != nil {
		panic(err)
	}

	opt.MaxActiveConns = int(conn.maxPoolConnection)
	opt.MinIdleConns = int(conn.minPoolConnection)
	opt.PoolTimeout = conn.waitPoolConnectionDuration
	opt.ConnMaxIdleTime = conn.maxIdleConnectionDuration
	opt.ConnMaxLifetime = conn.maxLifeTimeConnection
	opt.TLSConfig = conn.secureTLS

	rediscli := goredis.NewClusterClient(opt)

	cmd := rediscli.Ping(ctx)
	if cmd.Err() != nil {
		panic(cmd.Err())
	}

	return &RedisDBc{
		DB: rediscli,
	}
}

func SetRedisClusterURIConnection(uri string) OptionRedisCluster {
	return func(o *optionRedisCluster) {
		o.uri = uri
	}
}

func SetRedisClusterServiceName(serviceName string) OptionRedisCluster {
	return func(o *optionRedisCluster) {
		o.serviceName = serviceName
	}
}

func SetRedisClusterMinPoolConnection(minPoolConnection uint) OptionRedisCluster {
	return func(o *optionRedisCluster) {
		o.minPoolConnection = minPoolConnection
	}
}

func SetRedisClusterMaxPoolConnection(maxPoolConnection uint) OptionRedisCluster {
	return func(o *optionRedisCluster) {
		o.maxPoolConnection = maxPoolConnection
	}
}

func SetRedisClusterMaxIdleConnectionDuration(maxIdleConnectionDuration time.Duration) OptionRedisCluster {
	return func(o *optionRedisCluster) {
		o.maxIdleConnectionDuration = maxIdleConnectionDuration
	}
}

func SetRedisClusterWaitPoolConnectionDuration(waitPoolConnectionDuration time.Duration) OptionRedisCluster {
	return func(o *optionRedisCluster) {
		o.waitPoolConnectionDuration = waitPoolConnectionDuration
	}
}

func SetRedisClusterMaxLifeTimeConnection(maxLifeTimeConnection time.Duration) OptionRedisCluster {
	return func(o *optionRedisCluster) {
		o.maxLifeTimeConnection = maxLifeTimeConnection
	}
}

func SetRedisClusterSecureTLS(secureTLS *tls.Config) OptionRedisCluster {
	return func(o *optionRedisCluster) {
		o.secureTLS = secureTLS
	}
}

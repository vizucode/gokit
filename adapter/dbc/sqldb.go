package dbc

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"gitlab.com/tixia-backend/gokit/utils/env"
)

type OptionSQLDB func(*optionSqlDB)

type Driver string

const (
	Postgres Driver = "postgres"
	MySQL    Driver = "mysql"
)

type optionSqlDB struct {
	uri               string
	serviceName       string
	databaseName      string
	driver            Driver
	minPoolConnection uint
	maxPoolConnection uint
	maxConnectionIdle time.Duration
}

func defaultConnection() optionSqlDB {
	return optionSqlDB{
		uri:               env.GetString("DB_SQL_URI"),
		serviceName:       env.GetString("SERVICE_NAME"),
		databaseName:      env.GetString("DB_SQL_NAME"),
		minPoolConnection: 1,
		maxPoolConnection: 100,
		maxConnectionIdle: time.Minute * 1,
	}
}

func NewSqlConnection(options ...OptionSQLDB) *SqlDBC {
	var err error

	o := defaultConnection()
	for _, option := range options {
		option(&o)
	}

	dbc, err := sql.Open(string(o.driver), o.uri)
	if err != nil {
		panic(err)
	}

	dbc.SetConnMaxIdleTime(o.maxConnectionIdle)
	dbc.SetMaxIdleConns(int(o.minPoolConnection))
	dbc.SetMaxOpenConns(int(o.maxPoolConnection))

	err = dbc.Ping()
	if err != nil {
		log.Fatal("failed to connect to database")
	}

	return &SqlDBC{
		DB: dbc,
	}
}

func SetSqlURIConnection(uri string) OptionSQLDB {
	return func(o *optionSqlDB) {
		o.uri = uri
	}
}

func SetSqlDriver(driver Driver) OptionSQLDB {
	return func(o *optionSqlDB) {
		o.driver = driver
	}
}

func SetSqlServiceName(serviceName string) OptionSQLDB {
	return func(o *optionSqlDB) {
		o.serviceName = serviceName
	}
}

func SetSqlDatabaseName(databaseName string) OptionSQLDB {
	return func(o *optionSqlDB) {
		o.databaseName = databaseName
	}
}

func SetSqlMinPoolConnection(minPoolConnection uint) OptionSQLDB {
	return func(o *optionSqlDB) {
		o.minPoolConnection = minPoolConnection
	}
}

func SetSqlMaxPoolConnection(maxPoolConnection uint) OptionSQLDB {
	return func(o *optionSqlDB) {
		o.maxPoolConnection = maxPoolConnection
	}
}

func SetSqlMaxConnectionIdle(maxConnectionIdle time.Duration) OptionSQLDB {
	return func(o *optionSqlDB) {
		o.maxConnectionIdle = maxConnectionIdle
	}
}

package dbc

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/vizucode/gokit/utils/constant"
	"github.com/vizucode/gokit/utils/env"
)

type OptionSQLDB func(*optionSqlDB)

type optionSqlDB struct {
	uri               string
	serviceName       string
	databaseName      string
	driver            constant.Driver
	minPoolConnection uint
	maxPoolConnection uint
	maxConnectionIdle time.Duration
}

func defaultSqlDbConnection() optionSqlDB {
	return optionSqlDB{
		uri:               env.GetString("DB_SQL_URI"),
		serviceName:       env.GetString("SERVICE_NAME"),
		databaseName:      env.GetString("DB_SQL_NAME"),
		minPoolConnection: 1,
		maxPoolConnection: 100,
		maxConnectionIdle: time.Minute * 1,
	}
}

func NewSqlConnection(options ...OptionSQLDB) *SqlDBc {
	var err error

	o := defaultSqlDbConnection()
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

	return &SqlDBc{
		DB: dbc,
	}
}

func SetSqlURIConnection(uri string) OptionSQLDB {
	return func(o *optionSqlDB) {
		o.uri = uri
	}
}

func SetSqlDriver(driver constant.Driver) OptionSQLDB {
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

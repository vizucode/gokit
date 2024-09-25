package dbc

import (
	"time"

	"github.com/vizucode/gokit/utils/constant"
	"github.com/vizucode/gokit/utils/env"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type OptionsGormDB func(o *optionGormDB)

type optionGormDB struct {
	uri               string
	serviceName       string
	databaseName      string
	minPoolConnection uint
	maxPoolConnection uint
	skipTransaction   bool
	driver            constant.Driver
	maxIdleConnection time.Duration
}

func defaultGormConnection() optionGormDB {
	return optionGormDB{
		uri:               env.GetString("DB_GORM_URI"),
		serviceName:       env.GetString("SERVICE_NAME"),
		databaseName:      env.GetString("DB_GORM_NAME"),
		driver:            constant.Postgres,
		skipTransaction:   false,
		minPoolConnection: 1,
		maxPoolConnection: 100,
		maxIdleConnection: time.Minute * 1,
	}
}

func dialector(driver constant.Driver, dsn string) gorm.Dialector {
	switch driver {
	case constant.MySQL:
		return mysql.Open(dsn)
	default:
		return postgres.Open(dsn)
	}
}

func NewGormConnection(options ...OptionsGormDB) *GormDBC {
	var err error

	conn := defaultGormConnection()
	for _, o := range options {
		o(&conn)
	}

	gormDB, err := gorm.Open(dialector(conn.driver, conn.uri), &gorm.Config{
		SkipDefaultTransaction: conn.skipTransaction,
	})
	if err != nil {
		panic(err)
	}

	db, err := gormDB.DB()
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic("failed to connect to database")
	}

	db.SetMaxIdleConns(int(conn.minPoolConnection))
	db.SetMaxOpenConns(int(conn.maxPoolConnection))
	db.SetConnMaxIdleTime(conn.maxIdleConnection)

	return &GormDBC{
		DB: gormDB,
	}
}

func SetGormURIConnection(uri string) OptionsGormDB {
	return func(o *optionGormDB) {
		o.uri = uri
	}
}

func SetGormDriver(driver constant.Driver) OptionsGormDB {
	return func(o *optionGormDB) {
		o.driver = driver
	}
}

func SetGormServiceName(serviceName string) OptionsGormDB {
	return func(o *optionGormDB) {
		o.serviceName = serviceName
	}
}

func SetGormDatabaseName(databaseName string) OptionsGormDB {
	return func(o *optionGormDB) {
		o.databaseName = databaseName
	}
}

func SetGormMinPoolConnection(minPoolConnection uint) OptionsGormDB {
	return func(o *optionGormDB) {
		o.minPoolConnection = minPoolConnection
	}
}

func SetGormMaxPoolConnection(maxPoolConnection uint) OptionsGormDB {
	return func(o *optionGormDB) {
		o.maxPoolConnection = maxPoolConnection
	}
}

func SetGormSkipTransaction(skipTransaction bool) OptionsGormDB {
	return func(o *optionGormDB) {
		o.skipTransaction = skipTransaction
	}
}

func SetGormMaxIdleConnection(maxIdleConnection time.Duration) OptionsGormDB {
	return func(o *optionGormDB) {
		o.maxIdleConnection = maxIdleConnection
	}
}

package dbc

import (
	"database/sql"

	"gorm.io/gorm"
)

type SqlDBC struct {
	DB *sql.DB
}

type GormDBC struct {
	DB *gorm.DB
}

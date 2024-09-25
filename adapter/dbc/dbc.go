package dbc

import (
	"database/sql"

	"gorm.io/gorm"
)

// sqlDBC is instance for database/sql connection
type SqlDBC struct {
	DB *sql.DB
}

// Gorm is instance for gorm connection
type GormDBC struct {
	DB *gorm.DB
}

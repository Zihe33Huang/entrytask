package tcp

import (
	"database/sql"
)

var (
	db *sql.DB
)

func SetDB(d *sql.DB) {
	db = d
}

func GetDB() *sql.DB {
	return db
}

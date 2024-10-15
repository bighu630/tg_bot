package dao

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DbPath = "./quotations.db"

var db *sql.DB

func init() {
	ldb, err := sql.Open("sqlite3", DbPath)
	if err != nil {
		panic(err)
	}
	db = ldb
}

func GetDB() *sql.DB {
	return db
}

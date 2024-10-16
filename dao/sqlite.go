package dao

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const DbPath = "./quotations.db"

var db *sql.DB

func Init(path string) {
	if path == "" {
		path = DbPath
	}
	ldb, err := sql.Open("sqlite3", path)
	if err != nil {
		panic(err)
	}
	db = ldb
}

func GetDB() *sql.DB {
	return db
}

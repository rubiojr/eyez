package db

import "database/sql"

const (
	DefaultCaptureCollection = `requests`
	DefaultDatabase          = `records.db`
)

var db *sql.DB

const collectionCreateSQL = `CREATE TABLE IF NOT EXISTS "` + DefaultCaptureCollection + `" (
	"id" INTEGER PRIMARY KEY,
	"uuid" VARCHAR(36) NOT NULL,
	"method" VARCHAR(10),
	"status" INTEGER,
	"url" TEXT,
	"path" TEXT,
	"headers" TEXT,
	"body" BLOB,
	"date_start" DATETIME,
	"date_end" DATETIME,
	"time_taken" INTEGER
)`

func InitDB() error {
	var err error
	db, err = sql.Open("sqlite3", DefaultDatabase)
	if err != nil {
		return err
	}

	_, err = db.Exec(collectionCreateSQL)
	return err
}

func Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.Exec(query, args...)
}

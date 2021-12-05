package db

import "database/sql"

const (
	DefaultCaptureCollection = `requests`
	DefaultDatabase          = `records.db`
)

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

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", DefaultDatabase)
	if err != nil {
		return db, err
	}

	_, err = db.Exec(collectionCreateSQL)
	return db, err
}

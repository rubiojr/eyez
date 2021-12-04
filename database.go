package main

import "database/sql"

const (
	defaultCaptureCollection = `requests`
	defaultDatabase          = `records.db`
)

const collectionCreateSQL = `CREATE TABLE IF NOT EXISTS "` + defaultCaptureCollection + `" (
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

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", defaultDatabase)
	if err != nil {
		return db, err
	}
	defer db.Close()

	_, err = db.Exec(collectionCreateSQL)
	return db, err
}

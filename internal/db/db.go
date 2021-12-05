package db

import (
	"database/sql"
	"time"

	//_ "modernc.org/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

const (
	DefaultCaptureCollection = `requests`
	DefaultDatabase          = `records.db`
)

var db *sql.DB

type Record struct {
	ID        int
	UUID      string
	Method    string
	Status    int
	URL       string
	Path      string
	Headers   string
	Body      []byte
	DateStart *time.Time
	DateEnd   *time.Time
	TimeTaken int
}

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

func InitDB(path string) error {
	var err error
	db, err = sql.Open("sqlite3", path)
	if err != nil {
		return err
	}

	_, err = db.Exec(collectionCreateSQL)
	return err
}

func InitRODB(path string) error {
	var err error
	db, err = sql.Open("sqlite3", "file:"+path+"?mode=ro")
	return err
}

func Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.Exec(query, args...)
}

func Each(visitor func(rec *Record) error) error {
	rows, err := db.Query("SELECT id, uuid, method, status, url, path, headers, body, date_start, date_end FROM " + DefaultCaptureCollection + " ORDER BY date_end DESC")
	if err != nil {
		return err
	}

	for rows.Next() {
		rec := Record{}
		err := rows.Scan(&rec.ID, &rec.UUID, &rec.Method, &rec.Status, &rec.URL, &rec.Path, &rec.Headers, &rec.Body, &rec.DateStart, &rec.DateEnd)
		if err != nil {
			return err
		}
		if err := visitor(&rec); err != nil {
			return err
		}
	}

	return nil
}

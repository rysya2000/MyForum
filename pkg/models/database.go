package models

import (
	"database/sql"
	"io/ioutil"
)

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	//	Setup
	query, err := ioutil.ReadFile("./pkg/models/sqlite/setup.sql")
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(string(query)); err != nil {
		return nil, err
	}
	return db, nil
}

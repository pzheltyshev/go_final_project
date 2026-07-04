package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type DBConnection struct {
	DB *sql.DB
}

var Database DBConnection

const Schema string = `
CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(128) NOT NULL,
	comment TEXT,
	repeat VARCHAR(128)
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

func Init(dbPath string) error {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	Database = DBConnection{DB: db}

	_, err = Database.DB.Exec(Schema)
	return err
}

func GetDBConnection() DBConnection {
	return Database
}

package db

import (
	"database/sql"
	"fmt"
	"strconv"

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

func GetTask(id string) (*Task, error) {

	task := Task{}

	idTask, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	row := GetDBConnection().DB.QueryRow("SELECT date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", idTask))

	err = row.Scan(&task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func UpdateTask(task *Task) error {

	idTask, err := strconv.Atoi(task.ID)
	if err != nil {
		return err
	}

	query := "UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id"
	res, err := GetDBConnection().DB.Exec(query, sql.Named("date", task.Date), sql.Named("title", task.Title), sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat), sql.Named("id", idTask))
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

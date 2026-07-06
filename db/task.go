package db

import (
	"database/sql"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64

	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)"
	result, err := GetDBConnection().DB.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))

	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return id, nil
}

func Tasks(limit int) ([]*Task, error) {

	var tasks []*Task

	tasks = make([]*Task, 0)

	rows, err := GetDBConnection().DB.Query("SELECT id, date, title, comment, repeat FROM scheduler LIMIT :limit", sql.Named("limit", limit))

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

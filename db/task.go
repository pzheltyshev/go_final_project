package db

import (
	"database/sql"
	"fmt"
	"strconv"
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

func DeleteTask(id string) error {

	idTask, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	query := "DELETE FROM scheduler WHERE id = :id"
	_, err = GetDBConnection().DB.Exec(query, sql.Named("id", idTask))

	if err != nil {
		return err
	}

	return nil

}

func UpdateDate(nextDate string, id string) error {

	idTask, err := strconv.Atoi(id)

	if err != nil {
		return err
	}

	query := "UPDATE scheduler SET date = :date WHERE id = :id"
	res, err := GetDBConnection().DB.Exec(query, sql.Named("date", nextDate), sql.Named("id", idTask))
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

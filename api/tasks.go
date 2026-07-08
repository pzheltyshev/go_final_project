package api

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/pavel/go-final-project/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {

	tasks, err := db.Tasks(50)
	if err != nil {
		writeErrorJson(w, "Couldn't get tasks")
		return
	}

	writeJson(w, TasksResp{Tasks: tasks})
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	task, err := db.GetTask(id)

	if err != nil {
		writeErrorJson(w, "Couldn't fill task by id")
		return
	}

	task.ID = id

	writeJson(w, task)

}

func PutTaskHandler(w http.ResponseWriter, r *http.Request) {

	task := db.Task{}

	body, err := io.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		writeErrorJson(w, "Couldn't get updated task")
		return
	}

	err = json.Unmarshal(body, &task)
	if err != nil {
		writeErrorJson(w, "Invalid task format")
		return
	}

	if task.Title == "" {
		writeErrorJson(w, "Title is required")
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeErrorJson(w, err.Error())
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeErrorJson(w, "Couldn't update task")
		return
	}

	var emptyRes struct{}
	writeJson(w, emptyRes)
}

func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	task, err := db.GetTask(id)
	if err != nil {
		writeErrorJson(w, "Couldn't get task by id")
		return
	}

	if len(task.Repeat) == 0 {
		err = db.DeleteTask(id)
		if err != nil {
			writeErrorJson(w, "Couldn't delete task")
			return
		}

	} else {
		nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeErrorJson(w, "Couldn't get next date")
			return
		}

		err = db.UpdateDate(nextDate, id)
		if err != nil {
			writeErrorJson(w, "Couldn't update date")
			return
		}
	}

	var emptyRes struct{}
	writeJson(w, emptyRes)

}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	err := db.DeleteTask(id)
	if err != nil {
		writeErrorJson(w, "Couldn't delete task")
		return
	}

	var emptyRes struct{}
	writeJson(w, emptyRes)
}

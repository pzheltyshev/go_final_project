package api

import (
	"net/http"

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

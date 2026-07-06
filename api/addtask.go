package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/pavel/go-final-project/db"
)

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(DateFormat)
	}

	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return err
	}

	var nextDate string

	if len(task.Repeat) != 0 {
		nextDate, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if afterNow(now, t) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(DateFormat)
		} else {
			task.Date = nextDate
		}
	}

	return nil
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeErrorJson(w, err.Error())
		return
	}

	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		writeErrorJson(w, err.Error())
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

	id, err := db.AddTask(&task)
	if err != nil {
		writeErrorJson(w, err.Error())
		return
	}

	task.ID = strconv.FormatInt(id, 10)

	writeJson(w, struct {
		ID string `json:"id"`
	}{ID: task.ID})
}

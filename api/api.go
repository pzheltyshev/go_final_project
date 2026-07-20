package api

import (
	"encoding/json"
	"log"
	"net/http"
)

const DateFormat = "20060102"
const recordLimit = 50

func writeJson(w http.ResponseWriter, data any) {

	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)

	if err != nil {
		log.Printf("Failed to write response: %v", err)
		return
	}

}

func writeErrorJson(w http.ResponseWriter, error string, statusCode int) {

	resp, err := json.Marshal(struct {
		Error string `json:"error"`
	}{Error: error})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	_, err = w.Write(resp)

	if err != nil {
		log.Printf("Failed to write response: %v", err)
		return
	}

}

func Init(webPath string) {

	http.Handle("/", http.FileServer(http.Dir(webPath)))

	http.HandleFunc("GET /api/nextdate", NextDateHandler)

	http.HandleFunc("POST /api/task", AddTaskHandler)

	http.HandleFunc("GET /api/task", GetTaskHandler)

	http.HandleFunc("PUT /api/task", PutTaskHandler)

	http.HandleFunc("/api/tasks", tasksHandler)

	http.HandleFunc("POST /api/task/done", DoneTaskHandler)

	http.HandleFunc("DELETE /api/task", DeleteTaskHandler)

}

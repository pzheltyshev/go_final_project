package api

import (
	"encoding/json"
	"net/http"
	"path/filepath"
)

const DateFormat = "20060102"

func writeJson(w http.ResponseWriter, data any) {

	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

func writeErrorJson(w http.ResponseWriter, error string) {

	resp, err := json.Marshal(struct {
		Error string `json:"error"`
	}{Error: error})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(resp)

}

func Init(rootPath string) {

	webPath := filepath.Join(rootPath, "web")

	http.Handle("/", http.FileServer(http.Dir(webPath)))

	http.HandleFunc("/api/nextdate", NextDateHandler)

	http.HandleFunc("POST /api/task", AddTaskHandler)

	http.HandleFunc("/api/tasks", tasksHandler)

}

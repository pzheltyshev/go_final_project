package api

import (
	"net/http"
	"path/filepath"
)

const DateFormat = "20060102"

func Init(rootPath string) {

	webPath := filepath.Join(rootPath, "web")

	http.Handle("/", http.FileServer(http.Dir(webPath)))

	http.HandleFunc("/api/nextdate", NextDateHandler)

}

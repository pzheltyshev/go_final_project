package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pavel/go-final-project/db"
)

func main() {
	rootPath, err := os.Executable()
	if err != nil {
		log.Fatal("failed to get executable path:", err)
	}

	rootPath = filepath.Dir(rootPath)
	webPath := filepath.Join(rootPath, "web")
	dbPath := filepath.Join(rootPath, "scheduler.db")

	err = db.Init(dbPath)
	if err != nil {
		log.Fatal("failed to initialize database:", err)
	}

	http.Handle("/", http.FileServer(http.Dir(webPath)))

	err = http.ListenAndServe(":7540", nil)
	if err != nil {
		log.Fatal("failed to start server:", err)
	}

}

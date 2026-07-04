package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pavel/go-final-project/api"
	"github.com/pavel/go-final-project/db"
)

func Run() {
	rootPath, err := os.Executable()
	if err != nil {
		log.Fatal("failed to get executable path:", err)
	}

	rootPath = filepath.Dir(rootPath)
	dbPath := filepath.Join(rootPath, "scheduler.db")

	err = db.Init(dbPath)
	if err != nil {
		log.Fatal("failed to initialize database:", err)
	}

	api.Init(rootPath)

	err = http.ListenAndServe(":7540", nil)
	if err != nil {
		log.Fatal("failed to start server:", err)
	}

}

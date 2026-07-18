package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pavel/go-final-project/api"
	"github.com/pavel/go-final-project/db"
)

func Run() {
	rootPath, err := os.Executable()

	rootPath = filepath.Dir(rootPath)
	fmt.Println(rootPath)
	dirs, err := os.ReadDir(rootPath)

	for _, v := range dirs {

		fmt.Println(v.Name())

	}

	if err != nil {
		log.Fatal("failed to get executable path:", err)
	}

	rootPath = filepath.Dir(rootPath)

	dbPath := os.Getenv("TODO_DBFILE")
	if len(dbPath) == 0 {
		dbPath = filepath.Join(rootPath, "scheduler.db")
	}

	err = db.Init(dbPath)
	if err != nil {
		log.Fatal("failed to initialize database:", err)
	}

	defer db.GetDBConnection().DB.Close()

	api.Init(rootPath)

	port := os.Getenv("TODO_PORT")
	if len(port) == 0 {
		port = "7540"
	}

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("failed to start server:", err)
	}

}

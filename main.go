package main

import (
	"log"
	"os"

	"todo/pkg/db"
	"todo/pkg/server"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v\n", err)
	}
	defer db.DB.Close()

	server.Start()
}

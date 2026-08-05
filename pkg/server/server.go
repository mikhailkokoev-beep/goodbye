package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"todo/pkg/api"
)

func Start() {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	api.Init()

	http.Handle("/", http.FileServer(http.Dir("./web")))

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Запуск сервера на порту %s...\n", port)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}

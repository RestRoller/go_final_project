package server

import (
	"fmt"
	"net/http"
	"os"

	"go_final_project/pkg/api"
	"go_final_project/pkg/db"
)

func Run() error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	// Инициализируем БД (создаст таблицу если нет)
	if err := db.Init(dbFile); err != nil {
		return fmt.Errorf("ошибка инициализации БД: %w", err)
	}
	defer db.Close()

	// Регистрируем API обработчики
	api.Init()

	// Раздаём статические файлы
	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	addr := ":" + port
	fmt.Printf("Сервер запущен на http://localhost:%s\n", port)
	return http.ListenAndServe(addr, nil)
}

package main

import (
	"go_final_project/pkg/db"
	"log"
)

func initDB() {
	dbFile := "scheduler.db"
	if err := db.InitDB(dbFile); err != nil {
		log.Fatal("Ошибка инициализации БД:", err)
	}
	defer db.Close()
	log.Println("База данных успешно инициализирована")
}

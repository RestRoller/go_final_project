package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init(dbFile string) error {
	var err error

	log.Printf("Инициализация БД: %s", dbFile)

	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	// Создаём таблицу если её нет
	schema := `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT '',
		title VARCHAR(255) NOT NULL DEFAULT '',
		comment TEXT NOT NULL DEFAULT '',
		repeat VARCHAR(128) NOT NULL DEFAULT ''
	);
	CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
	`
	_, err = DB.Exec(schema)
	if err != nil {
		return err
	}

	log.Println("БД успешно инициализирована")
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

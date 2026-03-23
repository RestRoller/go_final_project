package db

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

// CreateTable создаёт таблицу scheduler в указанной БД
func CreateTable(db *sql.DB) error {
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
	_, err := db.Exec(schema)
	return err
}

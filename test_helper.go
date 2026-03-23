package main

import (
	"database/sql"
	"go_final_project/pkg/db"
)

func initTestDB(dbFile string) (*sql.DB, error) {
	return db.InitTest(dbFile)
}

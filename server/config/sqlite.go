package config

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitializeSQLite() (*gorm.DB, error) {

	dbPath := "./db/main.db"

	// Check if the database file exists
	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		log.Println("database file not found, creating...")
		// Create if the database file and directory
		err = os.MkdirAll("./db", os.ModePerm)
		if err != nil {
			panic(err)
		}
		file, err := os.Create(dbPath)
		if err != nil {
			panic(err)
		}
		file.Close()
	}

	// Create DB and connect
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Return the DB
	return db, nil
}

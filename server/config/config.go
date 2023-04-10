package config

import (
	"fmt"

	"gorm.io/gorm"
)

var (
	db *gorm.DB
)


func Init() error {
	var err error 

	// Inicialize SQLite
	db, err = InitializeSQLite()

	if err != nil {
		return fmt.Errorf("error initialize sqlite: %v", err)
	}

	return nil
}

func GetSQLite() *gorm.DB {
	return db
}
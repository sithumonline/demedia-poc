package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func Database(url string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	log.Print("database connected")
	return db
}

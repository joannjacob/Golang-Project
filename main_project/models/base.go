package models

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var db *gorm.DB

func init() {
	err := godotenv.Load("./config/local.env")
	if err != nil {
		log.WithFields(log.Fields{"method": "DB connection init()", "error": err}).Error("Error loading local.env file")
	}
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)
	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		log.WithFields(log.Fields{"method": "DB connection init()", "error": err}).Error("Error conecting to db")
	}
	db = conn
	db.Debug().AutoMigrate(&Account{}, &Product{})
}

func GetDB() *gorm.DB {
	return db
}

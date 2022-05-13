package models

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"index:idx_user,unique"`
	Hash     string `json:"hash"`
}

func ConnectDatabase() {
	// Load env vars
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	pgUser := os.Getenv("PG_USER")
	pgPass := os.Getenv("PG_PASS")
	pgDb := os.Getenv("PG_DB")
	pgHost := os.Getenv("PG_HOST")

	// Connect to DB
	dsn := "host=" + pgHost + " user=" + pgUser + " password=" + pgPass + " dbname=" + pgDb + " port=5432 sslmode=disable TimeZone=America/Chicago"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// Migrate the schema
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Error Migrating DB Schema")
		return
	}

	DB = db
}

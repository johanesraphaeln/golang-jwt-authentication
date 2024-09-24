package models

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	// "gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectRelationalDatabase() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("error loading .env file")
	}

	// Get database credentials from environment variables
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	databaseName := os.Getenv("DB_NAME")

	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, databaseName)
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, databaseName)

	// database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	database, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	database.AutoMigrate(&User{})

	DB = database
}
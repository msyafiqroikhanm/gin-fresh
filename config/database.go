package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupDatabase() (*gorm.DB, error) {
	// Load database credentials from environment variables
	var dbHost, dbPort, dbUser, dbPassword, dbName string
	if os.Getenv("PGHOST") == "production" {
		dbHost = os.Getenv("PROD_DB_HOST")
		dbPort = os.Getenv("PROD_DB_PORT")
		dbUser = os.Getenv("PROD_DB_USERNAME")
		dbPassword = os.Getenv("PROD_DB_PASSWORD")
		dbName = os.Getenv("PROD_DB_NAME")
	} else if os.Getenv("PGHOST") == "development" {
		dbHost = os.Getenv("DEV_DB_HOST")
		dbPort = os.Getenv("DEV_DB_PORT")
		dbUser = os.Getenv("DEV_DB_USERNAME")
		dbPassword = os.Getenv("DEV_DB_PASSWORD")
		dbName = os.Getenv("DEV_DB_NAME")
	} else {
		dbHost = os.Getenv("LOCAL_DB_HOST")
		dbPort = os.Getenv("LOCAL_DB_PORT")
		dbUser = os.Getenv("LOCAL_DB_USERNAME")
		dbPassword = os.Getenv("LOCAL_DB_PASSWORD")
		dbName = os.Getenv("LOCAL_DB_NAME")
	}

	// Construct database URL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Connect to PostgreSQL database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
		return nil, err
	}

	return db, nil
}

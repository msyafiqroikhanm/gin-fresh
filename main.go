package main

import (
	"jxb-eprocurement/config"
	"jxb-eprocurement/database"
	"jxb-eprocurement/handlers"
	"jxb-eprocurement/models"
	"jxb-eprocurement/routers"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// Load .env file
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	//use ../.env because main.go inside /cmd
	envPath := filepath.Join(pwd, "./.env")
	err = godotenv.Load(filepath.Join(envPath))
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Setup database connection
	db, err := config.SetupDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize DB in models
	models.InitDB(db)

	// Run migrations
	database.Migrate(db)

	// Setup router
	router := routers.SetupRouter(db)

	// Global error handler middleware
	router.Use(handlers.ErrorHandler)

	// Handle 404 Erorr Router
	router.Use(handlers.NotFoundHandler)

	// Force color output in gin
	gin.ForceConsoleColor()

	// Start server
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))

}

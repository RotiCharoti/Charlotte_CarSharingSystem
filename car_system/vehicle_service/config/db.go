package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/joho/godotenv"         // To load .env file
)

// DB is the global database connection pool
var DB *sql.DB

// ConnectDB initializes the connection to the MySQL database
func ConnectDB() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get database credentials from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Build the DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Initialize the database connection
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Verify the connection
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error verifying connection to the database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")
}

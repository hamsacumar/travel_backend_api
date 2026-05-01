package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	_ = godotenv.Load()

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPass == "" || dbName == "" {
		log.Fatal("Database environment variables are not set")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPass, dbName, "disable",
	)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to open DB connection:", err)
	}

	//connection pool tuning
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(10)

	err = DB.Ping()
	if err != nil {
		log.Fatal("Cannot connect to PostgreSQL:", err)
	}

	log.Println("✅ Connected to PostgreSQL successfully!")
}

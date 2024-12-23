package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func main() {
	// db connection details
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env file")
	}

	// read database configuration from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	// validate required environment variables
	if dbUser == "" || dbPassword == "" || dbName == "" || dbHost == "" || dbPort == "" {
		log.Fatal("Database environment variables are not set properly in the .env file")
	}

	// construct the Postgres connection string
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	// path to the SQL file
	initSQL := "db/init.sql"
	if len(os.Args) > 1 {
		initSQL = os.Args[1]
	}

	// read the SQL file
	schemaSQL, err := os.ReadFile(initSQL)
	if err != nil {
		log.Fatalf("failed to read .sql file: %v", err)
	}

	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatalf("error in connecting to the database: %v", err)
	}
	defer db.Close()

	// ping the database to verify the connection
	if err = db.Ping(); err != nil {
		log.Fatalf("could not ping the database: %v", err)
	}

	// execute the SQL statements
	start := time.Now()
	if _, err = db.Exec(string(schemaSQL)); err != nil {
		log.Fatalf("error in executing schema migration: %v", err)
	}

	fmt.Printf("DB migration completed in %v.\n", time.Since(start))
}

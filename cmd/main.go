package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"payment-gateway/db"
	"payment-gateway/internal/api"
	"payment-gateway/internal/kafka"
)

func init() {

}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env file")
	}

	// Read database configuration from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	kafkaURL := os.Getenv("KAFKA_BROKER_URL")

	// Validate required environment variables
	if dbUser == "" || dbPassword == "" || dbName == "" || dbHost == "" || dbPort == "" {
		log.Fatal("Database environment variables are not set properly in the .env file")
	}

	dbURL := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=disable"

	dbInst, err := db.New(dbURL)
	if err != nil {
		panic(err)
	}

	// Set up api endpoints
	router := api.New(dbInst)

	// init kafka
	kafkaInst := kafka.NewKafkaProducer(kafkaURL)

	router.SetupServices(kafkaInst)
	router.SetupRoutes()

	// Start the server on port 8080
	log.Println("Starting server on port 8090...")
	if err = http.ListenAndServe(":8090", router.Router); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}

}

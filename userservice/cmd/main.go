package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hari134/pratilipi/pkg/db"
	"github.com/hari134/pratilipi/pkg/kafka"
	"github.com/hari134/pratilipi/userservice/api"
	"github.com/hari134/pratilipi/userservice/middleware"
	"github.com/hari134/pratilipi/userservice/migrations"
	"github.com/hari134/pratilipi/userservice/producer"
)

func main() {
	// Load environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	serverPort := os.Getenv("SERVER_PORT")

	// Initialize the database using environment variables
	dbInstance := db.InitDB(db.Config{
		Host:     dbHost,
		Port:     stringToInt(dbPort, 5432), // default to 5432 if not set
		User:     dbUser,
		Password: dbPassword,
		DBName:   dbName,
		SSLMode:  "disable",
	})
	defer db.CloseDB(dbInstance)

	// Create Kafka configuration from environment variables
	kafkaConfig := kafka.NewKafkaConfig().
		SetBrokers(kafkaBrokers)

	// Initialize Kafka producer with KafkaConfig
	kafkaProducer := kafka.NewKafkaProducer(kafkaConfig)
	defer kafkaProducer.Close() // Close the Kafka producer gracefully

	// Initialize ProducerManager with KafkaProducer
	producerManager := producer.NewProducerManager(kafkaProducer)

	// Create API handlers
	userAPIHandler := &api.UserAPIHandler{
		DB:            dbInstance,
		KafkaProducer: producerManager,
	}

	authAPIHandler := &api.AuthAPIHandler{
		DB: dbInstance,
	}
	migrations.RunMigrations(dbInstance)
	// Set up HTTP router
	r := mux.NewRouter()
	r.HandleFunc("/login", authAPIHandler.LoginHandler).Methods("POST")
	r.HandleFunc("/validate-token", authAPIHandler.ValidateTokenHandler).Methods("POST")

	r.HandleFunc("/users/{userID}", userAPIHandler.GetUserByIdHandler).Methods("GET")
	r.HandleFunc("/users", userAPIHandler.GetUsersHandler).Methods("GET")
	r.HandleFunc("/create-user", userAPIHandler.CreateUserHandler).Methods("POST")
	r.Handle("/update-user", middleware.TokenValidationMiddleware(http.HandlerFunc(userAPIHandler.UpdateUserHandler))).Methods("PUT")

	// Start HTTP server
	log.Fatal(http.ListenAndServe(":"+serverPort, r))
}

// Utility function to convert string to int, with a default value fallback
func stringToInt(s string, defaultVal int) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return defaultVal
}

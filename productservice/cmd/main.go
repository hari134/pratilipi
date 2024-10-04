package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hari134/pratilipi/pkg/db"
	"github.com/hari134/pratilipi/pkg/kafka"
	"github.com/hari134/pratilipi/productservice/api"
	"github.com/hari134/pratilipi/productservice/consumer" // Import consumer package
	"github.com/hari134/pratilipi/productservice/migrations"
	"github.com/hari134/pratilipi/productservice/producer"
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
    defer db.CloseDB(dbInstance)
    migrations.RunMigrations(dbInstance)
    // Initialize Kafka producer
    kafkaConfig := kafka.NewKafkaConfig().
        SetBrokers(kafkaBrokers)

    kafkaProducer := kafka.NewKafkaProducer(kafkaConfig)

    // Initialize ProducerManager
    producerManager := producer.NewProducerManager(kafkaProducer)

    // Create API handlers
    productAPIHandler := &api.ProductAPIHandler{
        DB:       dbInstance,
        Producer: producerManager,
    }

    // Initialize Kafka consumer
    kafkaConsumerConfig := kafka.NewKafkaConfig().
        SetBrokers(kafkaBrokers)
    kafkaConsumerConfig.Topic = "order-placed"
    kafkaConsumer := kafka.NewKafkaConsumer(kafkaConsumerConfig)

    // Initialize ConsumerManager to listen for OrderPlaced events
    consumerManager := consumer.NewConsumerManager(kafkaConsumer, dbInstance)

    // Start listening to "Order Placed" events in a separate goroutine
    go consumerManager.StartOrderPlacedConsumer("order-placed")

    // Set up HTTP routes
    r := mux.NewRouter()
    r.HandleFunc("/products", productAPIHandler.CreateProductHandler).Methods("POST")
    r.HandleFunc("/products/{product_id}", productAPIHandler.UpdateProductHandler).Methods("PUT") // Update product
    r.HandleFunc("/products/{product_id}", productAPIHandler.DeleteProductHandler).Methods("DELETE") // Delete product
    r.HandleFunc("/products/{product_id}/inventory", productAPIHandler.UpdateInventoryHandler).Methods("PUT")

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
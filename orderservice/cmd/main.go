package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hari134/pratilipi/orderservice/api"
	"github.com/hari134/pratilipi/orderservice/consumer" // Import consumer package
	"github.com/hari134/pratilipi/pkg/messaging" // Import your message types
	"github.com/hari134/pratilipi/orderservice/migrations"
	"github.com/hari134/pratilipi/orderservice/producer"
	"github.com/hari134/pratilipi/pkg/db"
	"github.com/hari134/pratilipi/pkg/kafka"
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

	// Initialize Kafka producer
	kafkaConfig := kafka.NewKafkaConfig().
		SetBrokers(kafkaBrokers)

	kafkaProducer := kafka.NewKafkaProducer(kafkaConfig)

	// Initialize ProducerManager
	producerManager := producer.NewProducerManager(kafkaProducer)

	// Create API handlers
	orderAPIHandler := &api.OrderHandler{
		DB:       dbInstance,
		Producer: producerManager,
	}

	// Run DB migrations
	migrations.RunMigrations(dbInstance)

	// Initialize Kafka consumer
	kafkaConsumerConfig := kafka.NewKafkaConfig().
		SetBrokers("kafka:9092").
		SetGroupID("orderservice-group").
		SetGroupTopics("user-registered", "product-created", "order-placed") // Multiple topics

	kafkaConsumer := kafka.NewKafkaConsumer(kafkaConsumerConfig)

	// Initialize ConsumerManager
	consumerManager := consumer.NewConsumerManager(kafkaConsumer, dbInstance)

	// Register event types for dynamic handling
	kafkaConsumer.RegisterType("user-registered", &messaging.UserRegistered{})
	kafkaConsumer.RegisterType("product-created", &messaging.ProductCreated{})

	// Start listening to "User Registered" and "Product Created" events in separate goroutines

	go consumerManager.StartUserRegisteredConsumer("user-registered")
	go consumerManager.StartProductCreatedConsumer("product-created")

	// Set up HTTP routes
	r := mux.NewRouter()
	r.HandleFunc("/orders", orderAPIHandler.PlaceOrderHandler).Methods("POST")
	r.HandleFunc("/orders", orderAPIHandler.GetAllOrdersHandler).Methods("GET")            // Get all orders
	r.HandleFunc("/orders/{order_id}", orderAPIHandler.GetOrderByIDHandler).Methods("GET") // Get order by ID

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

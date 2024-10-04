package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hari134/pratilipi/orderservice/api"
	"github.com/hari134/pratilipi/orderservice/consumer" // Import consumer package
	"github.com/hari134/pratilipi/orderservice/producer"
	"github.com/hari134/pratilipi/pkg/db"
	"github.com/hari134/pratilipi/pkg/kafka"
	"github.com/hari134/pratilipi/orderservice/migrations"
)

func main() {
    // Initialize the database
    dbInstance := db.InitDB(db.Config{
        Host:     "localhost",
        Port:     5432,
        User:     "postgres",
        Password: "password",
        DBName:   "orderservice_db",
        SSLMode:  "disable",
    })
    defer db.CloseDB(dbInstance)

    // Initialize Kafka producer
    kafkaConfig := kafka.NewKafkaConfig().
        SetBrokers("localhost:9092")

    kafkaProducer := kafka.NewKafkaProducer(kafkaConfig)

    // Initialize ProducerManager
    producerManager := producer.NewProducerManager(kafkaProducer)

    // Create API handlers
    orderAPIHandler := &api.OrderHandler{
        DB:       dbInstance,
        Producer: producerManager,
    }
    migrations.RunMigrations(dbInstance)
    // Initialize Kafka consumer
    kafkaConsumerConfig := kafka.NewKafkaConfig().
        SetBrokers("localhost:9092").
        SetGroupID("order-service-group")

    kafkaConsumer := kafka.NewKafkaConsumer(kafkaConsumerConfig)

    // Initialize ConsumerManager to listen for "User Registered" and "Product Created" events
    consumerManager := consumer.NewConsumerManager(kafkaConsumer, dbInstance)

    // Start listening to "User Registered" and "Product Created" events in separate goroutines
    go consumerManager.StartUserRegisteredConsumer("user-registered")
    go consumerManager.StartProductCreatedConsumer("product-created")

    // Set up HTTP routes
    r := mux.NewRouter()
    r.HandleFunc("/orders", orderAPIHandler.PlaceOrderHandler).Methods("POST")
		r.HandleFunc("/orders", orderAPIHandler.GetAllOrdersHandler).Methods("GET")  // Get all orders
    r.HandleFunc("/orders/{order_id}", orderAPIHandler.GetOrderByIDHandler).Methods("GET")  // Get order by ID

    // Start HTTP server
    log.Fatal(http.ListenAndServe(":8083", r))
}

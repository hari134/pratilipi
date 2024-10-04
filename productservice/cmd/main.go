package main

import (
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/hari134/pratilipi/pkg/db"
    "github.com/hari134/pratilipi/pkg/kafka"
    "github.com/hari134/pratilipi/productservice/api"
    "github.com/hari134/pratilipi/productservice/migrations"
    "github.com/hari134/pratilipi/productservice/producer"
    "github.com/hari134/pratilipi/productservice/consumer" // Import consumer package
)

func main() {
    // Initialize the database
    dbInstance := db.InitDB(db.Config{
        Host:     "localhost",
        Port:     5432,
        User:     "postgres",
        Password: "password",
        DBName:   "productservice_db",
        SSLMode:  "disable",
    })
    defer db.CloseDB(dbInstance)
    migrations.RunMigrations(dbInstance)
    // Initialize Kafka producer
    kafkaConfig := kafka.NewKafkaConfig().
        SetBrokers("localhost:9092")

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
        SetBrokers("localhost:9092").
        SetGroupID("product-service-group")

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
    log.Fatal(http.ListenAndServe(":8082", r))
}

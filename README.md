# Microservices Architecture Project

This project consists of multiple microservices, each handling specific functionality within the system. The services include **User Service**, **Product Service**, **Order Service**, and **GraphQL Gateway**. Each microservice is containerized using Docker and communicates with other services via HTTP APIs.

## Table of Contents
- [Microservices Overview](#microservices-overview)
- [Features](#features)
- [Technologies Used](#technologies-used)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [Microservices List](#microservices-list)
- [REST API Endpoints](#rest-api-endpoints)
- [GraphQL API](#graphql-api)
- [Postman Collection](#postman-collection)
- [License](#license)

## Microservices Overview

- **User Service**: Handles user registration, authentication, and user management.
- **Product Service**: Manages product creation, updates, and listings.
- **Order Service**: Manages customer orders, including order creation and retrieval.
- **GraphQL Gateway**: A single API gateway that interacts with all microservices and exposes a unified GraphQL API to the clients.

## Features

- **User Management**: Register, authenticate, and fetch users.
- **Product Management**: Add, update, fetch, and list products.
- **Order Management**: Place orders, fetch orders by ID, and list all orders.
- **GraphQL API**: Centralized API to access all services via GraphQL.

## Technologies Used

- **Go (Golang)**: Backend for all microservices.
- **Kafka**: As a event stream for stream management.
- **PostgreSQL**: For persistent data storage.
- **GraphQL**: API for flexible data querying.
- **JWT**: For secure authentication and authorization.
- **Docker**: For containerization and managing microservices.

## Prerequisites

Ensure the following are installed:

- **Docker**: Version 19 or later.

## Getting Started

To start the entire system:

1. Clone the repository:

    ```bash
    git clone https://github.com/hari134/pratilipi.git
    cd pratilipi
    ```

2. Build and run the services:

    ```bash
    docker compose up --build
    ```

3. Create the kafka topics:

    ```bash
    docker exec -it kafka /bin/bash

    kafka-topics.sh --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic user-registered


    kafka-topics.sh --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic product-created

    kafka-topics.sh --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic order-placed


    kafka-topics.sh --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic inventory-updated
    ```
4. Access the GraphQL Playground at [http://localhost:8084](http://localhost:8084).

## Microservices List

- **User Service**: Accessible at port `8081`, responsible for user-related operations.
- **Product Service**: Accessible at port `8082`, manages product information.
- **Order Service**: Accessible at port `8083`, manages order processing.
- **GraphQL Gateway**: Accessible at port `8084`, consolidates all services under one API.

## REST API Endpoints

Below are the REST API endpoints exposed by each microservice:

- **User Service** (port `8081`):
    - `POST /users/register`: Register a new user.
    - `POST /users/login`: Authenticate a user.
    - `GET /users/{id}`: Fetch user details by ID.
    - `GET /users`: Fetch All users.

- **Product Service** (port `8082`):
    - `POST /products`: Create a new product.
    - `GET /products/{id}`: Fetch product details by ID.
    - `GET /products`: List all products.

- **Order Service** (port `8083`):
    - `POST /orders`: Create a new order.
    - `GET /orders/{id}`: Fetch order details by ID.
    - `GET /orders`: List all orders.

## GraphQL API

The GraphQL API supports:

- **User Queries**: Register, login, and fetch user data.
- **Product Queries**: Create, update, and fetch product information.
- **Order Queries**: Place orders, fetch orders by ID, and list all orders.

## Postman Collection

You can test the APIs with the following Postman collection:

[Postman Collection](https://www.postman.com/orbital-module-participant-42960309/workspace/pratilipi-hari/collection/6701938265f8ad9784cb5bd8?action=share&creator=38808772)

### Important Note:
- To interact with the GraphQL services, you need to **register a user** and then **log in** using the User Service to obtain a **JWT token**.
- **After logging in**, include the JWT token in the **Authorization header** of subsequent GraphQL requests.
- The `registerUser` mutation is the **only** GraphQL operation that does **not** require the `Authorization` header.

## License

This project is licensed under the [MIT License](LICENSE).

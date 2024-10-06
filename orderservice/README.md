# Order Service

The Order Service manages orders placed by users. It allows for order creation, fetching individual orders, and listing all orders. The service interacts with the GraphQL microservice to expose the order-related functionality via API.

## Table of Contents

- [Features](#features)
- [Technologies Used](#technologies-used)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [Authentication](#authentication)
- [API Endpoints](#api-endpoints)
- [Environment Variables](#environment-variables)
- [License](#license)

## Features

- **Order Creation**: Place an order with multiple items.
- **Order Retrieval**: Fetch an order by ID.
- **Order Listing**: Retrieve all orders placed by users.

## Technologies Used

- **Go (Golang)**: Backend implementation of the service.
- **PostgreSQL**: For storing orders and order items.
- **Docker**: For containerization of the service.
- **HTTP**: For communication with the Product and User services.

## Prerequisites

You will need the following tools installed:

- **Go**: Version 1.17 or above.
- **Docker**: Version 19 or above.
- **PostgreSQL**: For the order database.

## API Endpoints

- **POST /orders**: Place a new order.
- **GET /orders/{id}**: Fetch order details by ID.
- **GET /orders**: Retrieve all orders.

## License

This project is licensed under the MIT License.

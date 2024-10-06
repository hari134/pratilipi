# Product Service

The Product Service manages products available in the system. It allows for creating, updating, fetching, and listing products. The service interacts with the GraphQL microservice to expose product-related functionality via API.

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

- **Product Creation**: Add new products to the system.
- **Product Update**: Update existing product details.
- **Product Fetching**: Retrieve product details by ID.
- **Product Listing**: Fetch all available products.

## Technologies Used

- **Go (Golang)**: Backend implementation of the service.
- **PostgreSQL**: For product data storage.
- **Docker**: For containerization of the service.
- **HTTP**: For communication with other services.

## Prerequisites

You will need the following tools installed:

- **Go**: Version 1.17 or above.
- **Docker**: Version 19 or above.
- **PostgreSQL**: For the product database.

## API Endpoints

- **POST /products**: Create a new product.
- **GET /products/{id}**: Fetch product details by ID.
- **GET /products**: List all products.
- **PUT /products/{id}**: Update product details.

## License

This project is licensed under the MIT License.

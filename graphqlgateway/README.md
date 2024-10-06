# GraphQL Microservice

This project implements a GraphQL microservice that interacts with multiple backend services such as User Service, Product Service, and Order Service. The service exposes GraphQL queries and mutations, allowing clients to fetch and update data through a single API endpoint.

## Table of Contents

- [Features](#features)
- [Technologies Used](#technologies-used)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [Authentication](#authentication)
- [Queries and Mutations](#queries-and-mutations)
- [Environment Variables](#environment-variables)
- [License](#license)

## Features

- **User Management**: Register users, fetch user details, and list all users.
- **Product Management**: Create, fetch, and list products.
- **Order Management**: Place orders, fetch an order by ID, and list all orders.
- **Authentication**: JWT-based authentication with role validation for specific operations.

## Technologies Used

- **Go (Golang)**: For backend implementation.
- **GraphQL**: To expose a flexible querying API.
- **JWT**: For user authentication and authorization.
- **Docker**: For containerization of the service.
- **HTTP**: For communication with microservices (User, Product, Order services).

## Prerequisites

You will need the following tools installed:

- **Go**: Version 1.17 or above.
- **Docker**: Version 19 or above.
- **Microservices**: Ensure the User, Product, and Order services are running and accessible.

## Link to GraphQl collection
- https://www.postman.com/orbital-module-participant-42960309/workspace/pratilipi-hari/collection/6701938265f8ad9784cb5bd8?action=share&creator=38808772

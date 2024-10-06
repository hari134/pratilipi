# User Service

The User Service is responsible for managing user data, including registration, authentication, and profile management. It interacts with the GraphQL microservice to provide user-related functionality via API calls.

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

- **User Registration**: Allows users to register by providing their name, email, password, and phone number.
- **User Authentication**: Supports JWT-based authentication for login.
- **Profile Management**: Fetch and update user details.

## Technologies Used

- **Go (Golang)**: Backend implementation of the service.
- **PostgreSQL**: For user data storage.
- **JWT**: For user authentication.
- **Docker**: For containerization of the service.

## Prerequisites

You will need the following tools installed:

- **Go**: Version 1.17 or above.
- **Docker**: Version 19 or above.
- **PostgreSQL**: For the user database.

## API Endpoints

- **POST /register**: Register a new user.
- **POST /login**: Authenticate a user and return a JWT token.
- **GET /users/{id}**: Fetch user details by ID.
- **GET /users/**: Fetch all users.
- **POST /validate-token**: Validate JWT token and return user claims.

## License

This project is licensed under the MIT License.

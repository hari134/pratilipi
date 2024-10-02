# Start from the official Golang image
FROM golang:1.22-alpine

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go modules and dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application from the cmd directory
RUN go build -o userservice ./cmd

# Expose the port on which the application will run
EXPOSE 8080

# Run the application
CMD ["./userservice"]

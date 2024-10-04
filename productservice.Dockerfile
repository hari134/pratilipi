# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

# Install necessary dependencies
RUN apk --no-cache add git openssh

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go workspace (go.work, go.work.sum) and mod files
COPY go.work go.work.sum ./
COPY userservice/go.mod userservice/go.sum ./userservice/
COPY orderservice/go.mod orderservice/go.sum ./orderservice/
COPY productservice/go.mod productservice/go.sum ./productservice/
COPY graphqlgateway/go.mod graphqlgateway/go.sum ./graphqlgateway/

# Download all Go dependencies (using the workspace)
RUN go work sync && go mod download

# Copy the source code
COPY . .

# Build the Go binary for the userservice submodule
WORKDIR /app/productservice
RUN go build -o /app/productservicebin ./cmd

# Stage 2: Create a minimal runtime image
FROM alpine:latest

# Set the current working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/productservicebin .

# Ensure the binary has execution permissions
RUN chmod +x ./productservicebin

# Expose the port on which the application will run
EXPOSE 8080

# Run the precompiled binary
CMD ["./productservicebin"]

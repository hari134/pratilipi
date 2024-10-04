# Stage 1: Build the Go binaries
FROM golang:1.22 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Docker-specific go.work file and all service directories
COPY go.work.docker ./go.work
COPY userservice ./userservice
COPY orderservice ./orderservice
COPY productservice ./productservice
COPY graphqlgateway ./graphqlgateway

# Sync the Go workspace and download dependencies
RUN go work sync && go mod download

# Build the services
RUN go build -o /bin/userservice ./userservice/cmd
RUN go build -o /bin/orderservice ./orderservice/cmd
RUN go build -o /bin/productservice ./productservice/cmd
RUN go build -o /bin/graphqlgateway ./graphqlgateway/cmd

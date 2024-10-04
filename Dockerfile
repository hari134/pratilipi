# Stage 1: Build the Go binaries
FROM golang:1.22 AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Set the working directory inside the container
WORKDIR /app

# Copy the Go workspace and source code
COPY go.work.docker ./go.work
COPY userservice ./userservice
COPY orderservice ./orderservice
COPY productservice ./productservice

# Sync the Go workspace and download dependencies
RUN go work sync && go mod download

# Build the binaries for all services
RUN go build -o /app/bin/userservice ./userservice/cmd
RUN go build -o /app/bin/orderservice ./orderservice/cmd
RUN go build -o /app/bin/productservice ./productservice/cmd

# Stage 2: Create a minimal image for each service
FROM golang:1.22-alpine AS userservice
WORKDIR /app
COPY --from=builder /app/bin/userservice .
RUN chmod +x ./userservice

# Debug: List all files and directories before running the binary
RUN echo "Listing files in /app directory for userservice:" && ls -l /app

EXPOSE 8080
CMD ["./userservice"]

FROM golang:1.22-alpine AS orderservice
WORKDIR /app
COPY --from=builder /app/bin/orderservice .
RUN chmod +x ./orderservice

# Debug: List all files and directories before running the binary
RUN echo "Listing files in /app directory for orderservice:" && ls -l /app

EXPOSE 8080
CMD ["./orderservice"]

FROM golang:1.22-alpine AS productservice
WORKDIR /app
COPY --from=builder /app/bin/productservice .
RUN chmod +x ./productservice

# Debug: List all files and directories before running the binary
RUN echo "Listing files in /app directory for productservice:" && ls -l /app

EXPOSE 8080
CMD ["./productservice"]

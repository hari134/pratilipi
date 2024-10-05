# Stage 1: Build the Go binaries
FROM golang:1.22 AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Set the working directory inside the container
WORKDIR /app

# Copy the Go workspace and source code
COPY go.work ./go.work
COPY pkg ./pkg
COPY go.mod ./go.mod
COPY go.sum ./go.sum

COPY userservice ./userservice
COPY orderservice ./orderservice
COPY productservice ./productservice
COPY graphqlgateway ./graphqlgateway

# Sync the Go workspace and download dependencies
RUN go work sync && go mod download

# Build the binaries for all services
RUN go build -o /app/bin/userservice ./userservice/cmd
RUN go build -o /app/bin/orderservice ./orderservice/cmd
RUN go build -o /app/bin/productservice ./productservice/cmd
RUN go build -o /app/bin/graphqlgateway ./graphqlgateway/cmd

# Stage 2: Create a minimal image for each service
FROM golang:1.22-alpine AS userservice
WORKDIR /app
COPY --from=builder /app/bin/userservice .
RUN chmod +x ./userservice


EXPOSE 8080
CMD ["./userservice"]

FROM golang:1.22-alpine AS orderservice
WORKDIR /app
COPY --from=builder /app/bin/orderservice .
RUN chmod +x ./orderservice


EXPOSE 8080
CMD ["./orderservice"]

FROM golang:1.22-alpine AS productservice
WORKDIR /app
COPY --from=builder /app/bin/productservice .
RUN chmod +x ./productservice


EXPOSE 8080
CMD ["./productservice"]


FROM golang:1.22-alpine AS graphqlgateway
WORKDIR /app
COPY --from=builder /app/bin/graphqlgateway .
RUN chmod +x ./graphqlgateway


EXPOSE 8080
CMD ["./graphqlgateway"]

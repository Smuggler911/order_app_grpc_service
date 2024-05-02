# Builder stage
FROM golang:1.21-alpine as builder

RUN apk update && apk add --no-cache curl
# Create working directories
RUN mkdir /app/bin -p
RUN mkdir /bin/golang-migrate -p
# Download migrate app
RUN GOLANG_MIGRATE_VERSION=v4.15.1 && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/${GOLANG_MIGRATE_VERSION}/migrate.linux-amd64.tar.gz | \
    tar xvz migrate -C /bin/golang-migrate
# Set home directory
WORKDIR /app
# Copy go.mod
COPY go.mod go.sum /app/
# Download go dependencies
RUN go mod download
# Copy all local files
COPY . /app
# Build the Go application
RUN GOOS=linux go build -o /app/bin/app ./cmd/order




# MIGRATION (DEV)
FROM alpine:latest as dev-migration

# Install packages
RUN apk --no-cache  add ca-certificates

# Create home directory
WORKDIR /app
# Copy migration dir
COPY --from=builder /app/migrations/dev ./migrations
# Install migrate tool
COPY --from=builder /bin/golang-migrate /usr/local/bin



# Start a new stage using the Alpine Linux image
FROM alpine:latest as dev
# Install packages
RUN apk --no-cache add ca-certificates
# Create home directory
WORKDIR /app
# Copy the built executable from the builder stage
COPY --from=builder /app/bin/app /app/app
# Print the contents of /app after the COPY command
RUN ls
# Copy config file
COPY /config/local.yaml ./local.yaml
EXPOSE 8080
# Define the command to run the application
CMD ["./app","--config=./local.yaml"]
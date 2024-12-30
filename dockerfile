# Use the specific Go version
FROM golang:1.22.2-alpine

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Copy the .env file into the container
COPY configs/.env .env

# Build the application
RUN go build -o main .

# Expose the port
EXPOSE 8080

# Run the application
CMD ["./main"]

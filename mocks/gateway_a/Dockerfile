FROM golang:1.19-alpine

# Install necessary packages
RUN apk add --no-cache git netcat-openbsd gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum from the current directory
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go app
RUN go build -o /app/main .

# Command to run the executable 
CMD ["/app/main"]
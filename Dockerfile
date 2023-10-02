# Build stage
FROM golang:latest AS builder

# Set environment variables to disable CGO and ensure Go uses modules
ENV CGO_ENABLED=0
ENV GO111MODULE=on

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the application for an alpine-linux target
RUN go build -a -installsuffix cgo -o server cmd/server/main.go

# Final stage to produce the runtime image
FROM alpine:latest

# Update CA certificates
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy binary from the build stage
COPY --from=builder /app/server .

# Run the application
CMD ["./server"]

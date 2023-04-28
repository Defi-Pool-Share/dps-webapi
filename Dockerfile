# Build stage
FROM golang:1.20.3 as builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code to the working directory
COPY . .

# Build the application
RUN go build -o dps-webapi main.go

# Final stage
FROM gcr.io/distroless/base-debian10

# Copy the binary from the builder stage
COPY --from=builder /app/dps-webapi /dps-webapi

# Expose the port your application will run on
EXPOSE 8080

# Set the entrypoint
ENTRYPOINT ["/dps-webapi"]
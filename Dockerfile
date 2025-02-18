# Use the official Go image as a builder
FROM golang:1.24 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o app .

# Use a minimal runtime image
FROM gcr.io/distroless/base-debian12

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/app .

# Copy the embedded CSV file
COPY --from=builder /app/transformed_data.csv .

# Expose the port the application runs on
EXPOSE 3000

# Command to run the application
CMD ["/app/app"]

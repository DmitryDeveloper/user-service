# Stage 1: Build the Go application
FROM golang:1.21.1-alpine AS builder

# Adding git, bash and openssh to the image
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

RUN go mod vendor
# Build the Go app
RUN go build -o main .

# Stage 2: Create a minimal runtime image
FROM alpine:latest

# Add Maintainer Info
LABEL maintainer="Dzmitry Bahdanovich <dbogdanovich1993@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .
# Copy env file from root directory
COPY .env .

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["./main"]

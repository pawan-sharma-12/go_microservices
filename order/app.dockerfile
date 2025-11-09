# ----------------------
# Stage 1: Build the Go binary
# ----------------------
    FROM golang:1.21-alpine AS builder

    # Install required packages
    RUN apk --no-cache add gcc g++ make ca-certificates git
    
    # Set working directory
    WORKDIR /go/src/github.com/pawan-sharma-12/go_microservices
    
    # Copy Go modules
    COPY go.mod go.sum ./
    
    # Copy vendor directory if using vendoring
    COPY vendor vendor
    
    # Copy microservice code (and dependencies if needed)
    COPY order order
    COPY account account
    COPY catalog catalog
    
    # Build the binary
    RUN GO111MODULE=on go build -mod vendor -o /go/bin/app ./order/cmd/order
    
    # ----------------------
    # Stage 2: Minimal runtime image
    # ----------------------
    FROM alpine:latest
    
    # Install certificates for HTTPS
    RUN apk --no-cache add ca-certificates
    
    # Set working directory
    WORKDIR /usr/bin
    
    # Copy the built binary from builder
    COPY --from=builder /go/bin/app .
    
    # Expose the service port
    EXPOSE 8080
    
    # Command to run the microservice
    CMD ["app"]
    
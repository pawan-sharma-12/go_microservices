# ---------- Stage 1: Build ----------
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache gcc g++ make ca-certificates git

WORKDIR /go/src/github.com/pawan-sharma-12/go_microservices

COPY go.mod go.sum ./
COPY vendor vendor
COPY account account

RUN GO111MODULE=on go build -mod=vendor -o /go/bin/account ./account/cmd/account

# ---------- Stage 2: Run ----------
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /go/bin/account .

EXPOSE 8080

CMD ["./account"]

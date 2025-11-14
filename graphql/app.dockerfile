FROM golang:1.25-alpine AS builder

RUN apk add --no-cache gcc g++ make ca-certificates git

WORKDIR /go/src/github.com/pawan-sharma-12/go_microservices

COPY go.mod go.sum ./
COPY vendor vendor
COPY graphql graphql
COPY account account
COPY catalog catalog
COPY order order

RUN GO111MODULE=on go build -mod=vendor -o /go/bin/graphql-server ./graphql

FROM alpine:3.18

RUN apk add --no-cache tzdata ca-certificates

WORKDIR /app

COPY --from=builder /go/bin/graphql-server .

EXPOSE 8080

CMD ["./graphql-server"]

FROM golang:1.21-alpine AS builder
RUN apk --no-cache add gcc g++ make ca-certificates 
WORKDIR /go/src/github.com/pawan-sharma-12/go_microservices
COPY go.mod go.sum ./
COPY vendor vendor
COPY account account
RUN GO111MODULE=on go build -mod vendor -o /go/bin/app ./catalog/cmd/account

FROM alpine:latest
WORKDIR /usr/bin
COPY --from=builder /go/bin/app .
EXPOSE 8080
CMD ["app"]

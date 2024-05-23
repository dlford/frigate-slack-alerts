FROM golang:1.22.3-alpine3.19 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /dist/app .

FROM alpine:3.20.0
WORKDIR /dist
COPY --from=builder /dist/app .
CMD ["/dist/app"]
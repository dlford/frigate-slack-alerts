FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /dist/app .

FROM alpine:latest
WORKDIR /dist
COPY --from=builder /dist/app .
CMD ["/dist/app"]
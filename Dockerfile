FROM golang:1.22.5-alpine3.20 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /dist/app .

FROM alpine:3.20.2
WORKDIR /dist
COPY --from=builder /dist/app .
CMD ["/dist/app"]

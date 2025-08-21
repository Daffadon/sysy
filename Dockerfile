# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags "-s -w" -o app .

FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 2112
EXPOSE 5140

CMD ["./app"]
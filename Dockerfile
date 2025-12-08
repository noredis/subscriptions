FROM golang:1.25.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o subscriptions ./cmd/app

FROM alpine:3.20 AS runner

COPY --from=builder /app/subscriptions .

CMD ["./subscriptions"]

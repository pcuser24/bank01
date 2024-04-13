# Build stage
FROM golang:1.19.1-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o main
RUN apk add curl && (curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz)

# Run stage
FROM alpine:3.16
WORKDIR /app
COPY start.sh wait-for.sh ./
RUN chmod +x ./start.sh && chmod +x ./wait-for.sh
COPY --from=builder /app/migrate ./migrate
COPY db/migration ./migration
COPY .env.example app.env
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]

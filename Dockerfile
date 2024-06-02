# Build stage
FROM golang:1.22.3-alpine3.20 AS builder
WORKDIR /app
COPY . .
ENV CGO_ENABLED=1
RUN apk add build-base
RUN go build -o main cmd/api/main.go 
RUN go build -o migrator cmd/migrate/main.go 

# Run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/migrator .
COPY --from=builder /app/main .
COPY .env .

CMD [ "/app/migrator && /app/main" ]
EXPOSE 8081

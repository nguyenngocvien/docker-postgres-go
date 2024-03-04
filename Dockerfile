#Build stage
FROM golang:1.21-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main.go

#Run stage
FROM alpine3.18
WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080
CMD [ "/app/main" ]

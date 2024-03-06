#Build stage
FROM golang:1.21-alpine3.18 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

RUN go build -o main.go

#Install curl
RUN apk add curl

#Download and extract the migrate binary
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz

#Run stage, Create a minimal runtime image
FROM alpine:3.18

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/app .

#Copy downloaded migrate binary from the builder stage
COPY --from=builder /app/migate.linux-amd64 ./migrate

#Copy file
COPY app.env .
COPY start.sh .
COPY wait-for.sh .

#copy migration file
COPY ./db/migration ./db/migration

# Expose the port the application runs on
EXPOSE 8088

# Command to run the executable
CMD ["./app"]
ENTRYPOINT [ "./start.sh" , ]
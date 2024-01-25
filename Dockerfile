# syntax=docker/dockerfile:1
# Build stage
FROM golang:1.21-alpine as Builder
WORKDIR /students-api
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN GOOS=linux GOARCH=amd64 go build -o ./students-api cmd/main.go

# Final stage
FROM alpine:3.14
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=Builder /students-api/students-api .
EXPOSE 3000
CMD ["./students-api"]
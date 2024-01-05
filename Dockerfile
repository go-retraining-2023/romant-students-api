# syntax=docker/dockerfile:1

FROM golang:1.21-alpine
WORKDIR /students-api
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN GOOS=linux GOARCH=amd64 go build -o ./students-api cmd/main.go
EXPOSE 3000
CMD ["./students-api"]
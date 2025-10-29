FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go get -u github.com/gorilla/websocket
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /main .
FROM alpine:latest
COPY --from=builder /main /main
EXPOSE 8080
CMD ["/main"]
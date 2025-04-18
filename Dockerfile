
FROM golang:1.24.1 AS builder


WORKDIR /app


COPY go.mod ./
COPY main.go ./


RUN go mod download


RUN go build -o delete-by-query main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o delete-by-query main.go

FROM alpine:latest


WORKDIR /app


COPY --from=builder /app/delete-by-query .
RUN chmod +x /app/delete-by-query

ENTRYPOINT ["/app/delete-by-query"]

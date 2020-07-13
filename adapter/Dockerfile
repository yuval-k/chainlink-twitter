FROM golang:1.14.4 AS builder
WORKDIR /workdir/
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o adapter ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /workdir/adapter /usr/local/bin/adapter
ENTRYPOINT ["/usr/local/bin/adapter"]
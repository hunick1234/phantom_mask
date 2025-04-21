FROM golang:1.22.5-alpine as builder

RUN apk update && apk add --no-cache netcat-openbsd
# 設置工作目錄
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
COPY .env /app/.env

RUN go build -o phantomApp ./cmd/phantomApp/main.go
RUN go build -o seeder ./cmd/seeder/main.go

COPY init.sh /app/init.sh
RUN chmod +x /app/init.sh

ENTRYPOINT ["/app/init.sh"]

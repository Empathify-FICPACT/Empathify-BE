# Stage 1 — build
FROM golang:1.25-alpine AS builder

WORKDIR /app

# install dependencies
RUN apk add --no-cache git

RUN go install github.com/swaggo/swag/cmd/swag@latest

# copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# copy source code
COPY . .

# generate swagger docs
RUN swag init -g cmd/api/main.go --output docs

# build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Stage 2 — run
FROM alpine:3.19

WORKDIR /app

# install ca-certificates untuk HTTPS (Google OAuth butuh ini)
RUN apk --no-cache add ca-certificates tzdata

# set timezone
ENV TZ=Asia/Jakarta

# copy binary dari builder
COPY --from=builder /app/main .

EXPOSE 8008

CMD ["./main"]
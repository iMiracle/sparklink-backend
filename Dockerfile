FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o sparklink-api ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/sparklink-api .
COPY --from=builder /app/.env.example ./

EXPOSE 8080

ENV PORT=8080
ENV TZ=UTC

CMD ["./sparklink-api"]
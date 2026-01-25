FROM golang:1.25.5-alpine AS builder

WORKDIR /src
RUN apk add --no-cache git gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
  go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

RUN go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN cp /go/bin/migrate /out/migrate


FROM alpine:3.20

RUN apk add --no-cache ca-certificates wget su-exec

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /out/api /app/api
COPY --from=builder /out/migrate /usr/local/bin/migrate

COPY db/migrations /app/db/migrations
COPY db/seed /app/db/seed

COPY docker/entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

RUN mkdir -p /data && chown -R app:app /data

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=3s --retries=5 \
  CMD wget -qO- http://localhost:8080/api/ping || exit 1

USER root
ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["/app/api"]


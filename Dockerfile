# syntax=docker/dockerfile:1

# STAGE 1: BUILD the application
FROM golang:1.24-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /bin/app ./cmd/api

# STAGE 2: DEPLOY the application
FROM alpine:latest AS build-release-stage

WORKDIR /

RUN apk --no-cache add ca-certificates

# Create non-root user for security
RUN addgroup -S appuser \
    && adduser -S -G appuser -H -s /sbin/nologin appuser


COPY --from=build-stage --chown=appuser:appuser /bin/app /bin/app

EXPOSE 8080

# Run as non-root user
USER appuser

ENTRYPOINT ["/bin/app"]

# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.21.5-bookworm AS golang

# Copy go.mod and go.sum separately from the rest of the code,
# so their cached layer is not invalidated when the code changes.
COPY go.mod go.sum /
RUN go mod download

COPY . /app
WORKDIR /app/cmd/reviewbot

RUN go mod download
RUN go build -o=myreviewbot

# Production stage
FROM golang:1.21.5-bookworm

RUN pwd
RUN ls -al
COPY --from=golang /app/cmd/reviewbot/myreviewbot /app/myreviewbot
COPY --from=golang /app/internal/database/migrations /app/internal/database/migrations

EXPOSE ${HTTP_PORT:-4444}
WORKDIR /app/

CMD ["./myreviewbot"]

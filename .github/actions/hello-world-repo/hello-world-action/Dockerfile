# syntax=docker/dockerfile:1

# --- build stage ---
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 go build -o /hello-world-action .

# --- final stage: just the binary, nothing else ---
FROM alpine:3.19
COPY --from=build /hello-world-action /hello-world-action
ENTRYPOINT ["/hello-world-action"]

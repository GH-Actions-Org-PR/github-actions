# syntax=docker/dockerfile:1
FROM golang:1.22-bookworm AS build

WORKDIR /src

# Cache module downloads separately from source changes.
COPY go.mod ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .

ARG VERSION=dev
ARG COMMIT=none
ARG BUILD_DATE=unknown

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-s -w \
      -X github.com/example/go-api/internal/version.Version=${VERSION} \
      -X github.com/example/go-api/internal/version.Commit=${COMMIT} \
      -X github.com/example/go-api/internal/version.BuildDate=${BUILD_DATE}" \
    -o /out/api ./cmd/api

# Distroless: no shell, no package manager, minimal attack surface.
FROM gcr.io/distroless/static-debian12:nonroot AS runtime

COPY --from=build /out/api /api

USER nonroot:nonroot
EXPOSE 8080

ENTRYPOINT ["/api"]

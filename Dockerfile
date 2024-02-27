# Build
FROM golang:1.22-bullseye AS build

WORKDIR /usr/src/app

COPY go.mod ./

RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go mod download

# Development
FROM build AS dev

RUN go install github.com/cosmtrek/air@latest && \
  go install github.com/go-delve/delve/cmd/dlv@latest

COPY . .

CMD ["air", "-c", ".air.toml"]

# Build production binaries
FROM build AS build-production

COPY . .

RUN go build \
  -ldflags="-linkmode external -extldflags -static" \
  -tags netgo \
  -o bin/go-tasks-manager ./cmd

# Production
FROM scratch

ENV GIN_MODE=release

WORKDIR /

COPY --from=build-production /usr/src/app/bin/go-tasks-manager go-tasks-manager

EXPOSE 8080

CMD ["/go-tasks-manager"]
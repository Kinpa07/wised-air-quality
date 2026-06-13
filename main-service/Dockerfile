# --- Build stage -----------------------------------------------------------
# Compiles the binary inside Docker so no local Go/C toolchain is required.
FROM golang:1.26-alpine AS builder

# CGO is required by the sqlite3 driver, so pull in a C toolchain.
RUN apk add --no-cache gcc musl-dev

WORKDIR /src

# Vendored dependencies are committed, so no module download is needed.
COPY go.mod go.sum ./
COPY vendor/ vendor/
COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -mod=vendor -o /service cmd/root.go

# --- Runtime stage ---------------------------------------------------------
FROM alpine:3.21.2

# musl is already present in alpine; add certs for outbound TLS if needed.
RUN apk add --no-cache ca-certificates

COPY config/config.yaml /config.yaml

# Districts GeoJSON read at startup; binary runs from / so assets/ resolves here.
COPY assets/ /assets/

COPY --from=builder /service /service

EXPOSE 8080

CMD [ "/service" ]

ARG GO_VERSION=1.27.1

FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates tzdata

# Keep dependency download cacheable independently from source changes.
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/products ./cmd

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=builder /out/products /app/products
COPY --from=builder /src/internal/connector/db/migrations /app/internal/connector/db/migrations

EXPOSE 8090
ENTRYPOINT ["/app/products"]
FROM golang:1.27.3 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/
RUN CGO_ENABLED=0 go build -o /out/slog-otelcol-playground ./cmd/slog-otelcol-playground

FROM alpine:3

RUN apk add --no-cache ca-certificates

COPY --from=builder /out/slog-otelcol-playground /usr/local/bin/slog-otelcol-playground

ENTRYPOINT ["/usr/local/bin/slog-otelcol-playground"]

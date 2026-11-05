FROM golang:1.27.3

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd/
COPY internal ./internal/

ENTRYPOINT ["go", "run", "./cmd/slog-otelcol-playground"]

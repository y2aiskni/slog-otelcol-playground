package util

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	slogltsv "github.com/ophum/slog-ltsv"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

func NewSlogLTSVHandler() slog.Handler {
	return slogltsv.NewLTSVHandler(os.Stdout, &slogltsv.HandlerOptions{
		Level:      slog.LevelDebug,
		TimeFormat: time.RFC3339Nano,
	})
}

func NewSlogOTLPHandler(ctx context.Context, serviceName string, endpointURL string) (slog.Handler, func(context.Context) error, error) {
	exporter, err := otlploggrpc.New(ctx, otlploggrpc.WithEndpointURL(endpointURL))
	if err != nil {
		return nil, nil, fmt.Errorf("create OTLP exporter: %w", err)
	}
	provider := sdklog.NewLoggerProvider(
		sdklog.WithResource(resource.NewWithAttributes("", attribute.String("service.name", serviceName))),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
	)
	return otelslog.NewHandler(serviceName, otelslog.WithLoggerProvider(provider)), provider.Shutdown, nil
}

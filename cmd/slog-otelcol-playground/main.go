package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/y2aiskni/go-slog-otelcol-playground/internal/util"
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// 標準出力(LTSV形式)とOTLPの両方にログを出力する
	stdoutLogger := util.NewSlogLTSVHandler()
	otlpLogger, shutdown, err := util.NewSlogOTLPHandler(ctx, "slog-otelcol-playground", "http://otelcol:4317")
	if err != nil {
		return fmt.Errorf("create OTLP handler: %w", err)
	}
	slog.SetDefault(slog.New(slog.NewMultiHandler(stdoutLogger, otlpLogger)))

	// 雑にログを出力する
	slog.DebugContext(ctx, "debug log", slog.Bool("enabled", true))
	slog.InfoContext(ctx, "hello from slog", slog.Int("count", 1))
	slog.With(slog.String("kind", "http")).InfoContext(ctx, "request completed",
		slog.String("http_request_method", "GET"),
		slog.String("http_request_url", "https://example.com/api/v1/resource"),
		slog.Int("http_request_size", 123),
		slog.Int("http_status_code", 200),
		slog.Int("http_response_size", 456),
		slog.String("http_user_agent", "Mozilla/5.0"),
		slog.String("http_remote_ip", "203.0.113.1"),
		slog.Int("http_remote_port", 12345),
		slog.String("http_referer", "https://example.com/previous"),
		slog.Int64("http_latency_ns", 123456789),
		slog.Bool("http_cache_lookup", true),
		slog.Bool("http_cache_hit", true),
		slog.Int("http_cache_fill_bytes", 0),
		slog.String("http_protocol", "HTTP/1.1"),
		slog.String("http_scheme", "https"),
		slog.String("http_user", "alice"),
	)
	slog.WarnContext(ctx, "warning log", slog.String("message", "tab\tand newline\nare escaped"))
	slog.ErrorContext(ctx, "sample error log", slog.Any("error", errors.New("example error")))

	// 終了時に溜まったOTLPログを送り込む
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("flush OTLP logs: %w", err)
	}
	return nil
}

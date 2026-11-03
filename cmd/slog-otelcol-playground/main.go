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
	slog.Default().WithGroup("request").InfoContext(ctx, "request completed", slog.String("method", "GET"), slog.Int("status", 200))
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

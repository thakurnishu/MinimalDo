package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.32.0"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	tracer trace.Tracer
	logger *slog.Logger
)

type MultiHandler struct {
	handlers []slog.Handler
}

func InitTelemetry(cfg *Config) (func(), error) {
	ctx := context.Background()

	// OTEL resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			attribute.String("library.language", "go"),
		),
	)
	errorLog("Failed to create otlp resource", err)

	// Initialize tracing
	traceCleanup, err := initTracing(ctx, res, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tracing: %w", err)
	}

	// Initialize logging
	logCleanup, err := initLogging(ctx, res, cfg)
	if err != nil {
		traceCleanup()
		return nil, fmt.Errorf("failed to initialize logging: %w", err)
	}

	return func() {
		logCleanup()
		traceCleanup()
	}, nil
}

func initTracing(ctx context.Context, res *resource.Resource, cfg *Config) (func(), error) {
	// OTLP GRPC trace exporter
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(cfg.SignozEndpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create trace exportor: %w", err)
	}

	// Trace provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Set global trace provider
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Tracer
	tracer = otel.Tracer(cfg.ServiceName)

	return func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tp.Shutdown(shutdownCtx); err != nil {
			slog.Error("Error shutting down tracer provider", "error", err)
		}
	}, nil
}

func initLogging(ctx context.Context, res *resource.Resource, cfg *Config) (func(), error) {
	// OTLP GRPC log exporter
	logExporter, err := otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(cfg.SignozEndpoint),
		otlploggrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create log exporter: %w", err)
	}

	// Log processor
	logProcessor := sdklog.NewBatchProcessor(logExporter)

	// Logger Provider
	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(logProcessor),
		sdklog.WithResource(res),
	)

	// Set global logger provider
	global.SetLoggerProvider(loggerProvider)	

	// Handlers
	var handlers []slog.Handler
	var logHandler slog.Handler

	// OpenTelemetry slog handler
	otelHander := otelslog.NewHandler(cfg.ServiceName)
	handlers = append(handlers, otelHander)

	// check if console log is enable
	if cfg.EnableConsoleLog {
		consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
		handlers = append(handlers, consoleHandler)
	}

	if len(handlers) == 1 {
		logHandler = handlers[0]
	}	else {
		logHandler = NewMultiHandler(handlers...)
	}

	// Set global logger
	logger = slog.New(logHandler)
	slog.SetDefault(logger)

	return func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := loggerProvider.Shutdown(shutdownCtx); err != nil {
			slog.Error("Error shutting down logger provider", "error", err)
		}
	}, nil
}

// Implements slog.Handler interface methods
func NewMultiHandler(handlers ...slog.Handler) *MultiHandler {
	return &MultiHandler{handlers: handlers}
}

func (h *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *MultiHandler) Handle(ctx context.Context, record slog.Record) error {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, record.Level) {
			if err := handler.Handle(ctx, record); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithAttrs(attrs)
	}

  return &MultiHandler{handlers: newHandlers}
}

func (h *MultiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithGroup(name)
	}
	return &MultiHandler{handlers: newHandlers}
}

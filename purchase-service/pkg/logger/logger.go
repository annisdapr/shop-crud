// 

package logger

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	TraceID   string `json:"trace_id,omitempty"`
	SpanID    string `json:"span_id,omitempty"`
}

var stdLogger = log.New(os.Stdout, "", 0) 

func Info(ctx context.Context, message string) {
	logWithLevel(ctx, "INFO", message)
}

func Warn(ctx context.Context, message string) {
	logWithLevel(ctx, "WARN", message)
}

func Error(ctx context.Context, message string) {
	logWithLevel(ctx, "ERROR", message)
}

func logWithLevel(ctx context.Context, level string, message string) {
	spanCtx := trace.SpanContextFromContext(ctx)

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		Message:   message,
	}

	if spanCtx.HasTraceID() {
		entry.TraceID = spanCtx.TraceID().String()
	}
	if spanCtx.HasSpanID() {
		entry.SpanID = spanCtx.SpanID().String()
	}

	logLine, err := json.Marshal(entry)
	if err != nil {
		// Fallback if marshalling failed
		log.Printf("Error marshalling log: %v", err)
		log.Printf("[%s] [trace_id=%s span_id=%s] %s", level, entry.TraceID, entry.SpanID, message)
		return
	}

	stdLogger.Println(string(logLine))
}

func LogWithTrace(ctx context.Context, message string) {
	Info(ctx, message)
}
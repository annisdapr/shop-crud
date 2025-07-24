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

// Definisikan struktur log kustom
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	TraceID   string `json:"trace_id,omitempty"`
	SpanID    string `json:"span_id,omitempty"`
}

var stdLogger = log.New(os.Stdout, "", 0) // Hapus flag default agar output bersih

// LogWithTrace sekarang mencatat dalam format JSON
func LogWithTrace(ctx context.Context, message string) {
	spanCtx := trace.SpanContextFromContext(ctx)

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     "INFO",
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
		// Fallback jika marshalling gagal
		log.Printf("Error marshalling log: %v", err)
		log.Printf("[trace_id=%s span_id=%s] %s", entry.TraceID, entry.SpanID, message)
		return
	}

	stdLogger.Println(string(logLine))
}
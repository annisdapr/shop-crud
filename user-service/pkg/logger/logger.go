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

// --- FUNGSI HELPER BARU UNTUK SETIAP LEVEL ---

// Info mencatat pesan dengan level "INFO"
func Info(ctx context.Context, message string) {
	logWithLevel(ctx, "INFO", message)
}

// Warn mencatat pesan dengan level "WARN"
func Warn(ctx context.Context, message string) {
	logWithLevel(ctx, "WARN", message)
}

// Error mencatat pesan dengan level "ERROR"
func Error(ctx context.Context, message string) {
	logWithLevel(ctx, "ERROR", message)
}

// --- FUNGSI INTI YANG DIPERBARUI ---

// logWithLevel adalah fungsi inti yang sekarang menerima parameter level.
// Fungsi ini bersifat private karena kita akan memanggilnya melalui helper di atas.
func logWithLevel(ctx context.Context, level string, message string) {
	spanCtx := trace.SpanContextFromContext(ctx)

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level, // <-- Menggunakan level dari parameter
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
		log.Printf("[%s] [trace_id=%s span_id=%s] %s", level, entry.TraceID, entry.SpanID, message)
		return
	}

	stdLogger.Println(string(logLine))
}

// LogWithTrace sekarang menjadi alias untuk Info agar kompatibel dengan kode lama.
func LogWithTrace(ctx context.Context, message string) {
	Info(ctx, message)
}
// tracing.go
package tracing

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func InitTracerProvider(serviceName, endpoint string) (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	// Opsi exporter OTLP, akan kita sesuaikan di bawah
	exporterOptions := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
	}

	// Baca variabel lingkungan untuk menentukan backend
	backend := os.Getenv("TRACING_BACKEND")

	// Jika backend-nya adalah OpenObserve, tambahkan konfigurasi spesifiknya
	if backend == "openobserve" {
		token := os.Getenv("ZO_ROOT_USER_TOKEN")
		if token != "" {
			headers := map[string]string{
				"Authorization": "Basic " + token,
			}
			exporterOptions = append(exporterOptions, otlptracehttp.WithHeaders(headers))
		}
		exporterOptions = append(exporterOptions, otlptracehttp.WithURLPath("/api/default/v1/traces"))
	}
	// Jika backend adalah "tempo" atau tidak diatur, tidak perlu opsi tambahan.
	// Exporter akan menggunakan path default /v1/traces yang benar untuk Tempo.
	
	exporter, err := otlptracehttp.New(ctx, exporterOptions...)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	return tp, nil
}
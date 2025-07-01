package tracing

import (
   "context"
   "fmt"
   "os"

   "go.opentelemetry.io/otel"
   "go.opentelemetry.io/otel/attribute"
   "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
   "go.opentelemetry.io/otel/sdk/resource"
   sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// getEnvOrDefault returns the value of the environment variable named by key
// or fallback if the variable is empty.
func getEnvOrDefault(key, fallback string) string {
   if val := os.Getenv(key); val != "" {
       return val
   }
   return fallback
}

// InitTracerProvider sets up an OpenTelemetry TracerProvider with an OTLP HTTP exporter.
// serviceName: name of the service generating traces.
// collectorHost: host (with port) of the OTLP HTTP collector (no scheme).
// The exporter sends to /api/default/v1/traces and uses an Authorization header
// with token read from OTEL_AUTH_TOKEN environment variable (fallback empty).
// It also sets resource attributes for service.name and environment=development,
// and uses AlwaysSample sampler.
func InitTracerProvider(serviceName, collectorHost string) (*sdktrace.TracerProvider, error) {
   ctx := context.Background()
   token := getEnvOrDefault("OTEL_AUTH_TOKEN", "")
   headers := make(map[string]string)
   if token != "" {
       headers["Authorization"] = fmt.Sprintf("Bearer %s", token)
   }

   exp, err := otlptracehttp.New(ctx,
       otlptracehttp.WithEndpoint(collectorHost),
       otlptracehttp.WithURLPath("/api/default/v1/traces"),
       otlptracehttp.WithHeaders(headers),
       otlptracehttp.WithInsecure(),
   )
   if err != nil {
       return nil, err
   }

   res, err := resource.New(ctx,
       resource.WithAttributes(
           attribute.String("service.name", serviceName),
           attribute.String("environment", "development"),
       ),
   )
   if err != nil {
       return nil, err
   }

   tp := sdktrace.NewTracerProvider(
       sdktrace.WithSampler(sdktrace.AlwaysSample()),
       sdktrace.WithResource(res),
       sdktrace.WithBatcher(exp),
   )
   otel.SetTracerProvider(tp)
   return tp, nil
}
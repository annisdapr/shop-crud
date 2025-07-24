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

// getEnvOrDefault returns the value of the environment variable key,
// or fallback if the variable is not set.
func getEnvOrDefault(key, fallback string) string {
   if v := os.Getenv(key); v != "" {
       return v
   }
   return fallback
}

// InitTracerProvider initializes an OpenTelemetry TracerProvider with an OTLP HTTP exporter.
// serviceName identifies this service. collectorHost is the OTLP endpoint host (host:port).
func InitTracerProvider(serviceName, collectorHost string) (*sdktrace.TracerProvider, error) {
   ctx := context.Background()

   token := getEnvOrDefault("ZO_ROOT_USER_TOKEN", "")
   headers := map[string]string{}
   if token != "" {
       headers["Authorization"] = "Basic " + token
   }
   exporter, err := otlptracehttp.New(
       ctx,
       otlptracehttp.WithEndpoint(collectorHost),
       otlptracehttp.WithURLPath("/api/default/v1/traces"),
       otlptracehttp.WithHeaders(headers),
       otlptracehttp.WithInsecure(),
   )
   if err != nil {
       return nil, err
   }

   res, err := resource.New(
       ctx,
       resource.WithAttributes(
           semconv.ServiceNameKey.String(serviceName),
           semconv.DeploymentEnvironmentKey.String("development"),
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
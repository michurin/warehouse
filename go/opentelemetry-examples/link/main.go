package main

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func initTracer() func(context.Context) error {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint()) // write to stdout
	if err != nil {
		panic(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.Default()),
	)

	otel.SetTracerProvider(tp)

	return tp.Shutdown
}

func g(ctx context.Context) {
	tracer := otel.Tracer("demo-tracer")

	ctx, span := tracer.Start(ctx, "g")
	defer span.End()

	time.Sleep(1 * time.Second)
}

func f(ctx context.Context) {
	tracer := otel.Tracer("demo-tracer")

	ctx, span := tracer.Start(ctx, "f", trace.WithAttributes(attribute.String("CUSTOM", "VALUE")))
	defer span.End()

	g(ctx)
	time.Sleep(1 * time.Second)
}

func main() {
	ctx := context.Background()

	shutdown := initTracer()
	defer shutdown(ctx)

	tracer := otel.Tracer("demo-tracer")

	ctx, span := tracer.Start(ctx, "main")
	defer span.End()

	f(ctx)

	sctx := span.SpanContext()

	ctx, spanB := tracer.Start(
		context.Background(), // новый root context
		"span-B",
		trace.WithLinks(trace.Link{
			SpanContext: sctx,
		}),
	)
	defer spanB.End()

	g(ctx)

	time.Sleep(1 * time.Second)
}

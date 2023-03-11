package xtrace

import (
	"context"

	sdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	globalProvider = sdk.NewTracerProvider()
	globalTracer   = globalProvider.Tracer("Default")
)

func GlobalTracer() trace.Tracer {
	return globalTracer
}

func SetGlobalTracer(tracer trace.Tracer) {
	globalTracer = tracer
}

func SetGlobalProvider(opts ...sdk.TracerProviderOption) *sdk.TracerProvider {
	return sdk.NewTracerProvider(opts...)
}

func NewTracer(name string, opts ...trace.TracerOption) trace.Tracer {
	return globalProvider.Tracer(name, opts...)
}

func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return GlobalTracer().Start(ctx, name, opts...)
}

func ParseOrGenTraceID(ctx context.Context) string {
	_, span := GlobalTracer().Start(ctx, "ParseOrGenTraceID")
	return TraceID(span)
}

func TraceID(span trace.Span) string {
	return span.SpanContext().TraceID().String()
}

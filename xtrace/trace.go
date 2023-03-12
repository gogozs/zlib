package xtrace

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	sdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/status"
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

func Inject(ctx context.Context, p propagation.TextMapPropagator, carrier propagation.TextMapCarrier) {
	p.Inject(ctx, carrier)
}

func Extract(ctx context.Context, p propagation.TextMapPropagator, carrier propagation.TextMapCarrier) context.Context {
	return p.Extract(ctx, carrier)
}

func End(span trace.Span, err error) {
	defer span.End()
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			span.SetStatus(codes.Error, s.Message())
		} else {
			span.SetStatus(codes.Error, err.Error())
		}
		return
	}
	span.SetStatus(codes.Ok, "success")
}

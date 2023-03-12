package xtrace

import (
	"context"

	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/metadata"
)

func InjectGrpcTraceInfo(ctx context.Context) metadata.MD {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.Pairs()
	}
	Inject(ctx, propagation.TraceContext{}, NewGrpcMetaData(md))
	return md
}

func ExtractGrpcTraceInfo(ctx context.Context) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.Pairs()
	}
	return Extract(ctx, propagation.TraceContext{}, NewGrpcMetaData(md))
}

type GrpcMetaData struct {
	md metadata.MD
}

var _ propagation.TextMapCarrier = (*GrpcMetaData)(nil)

func NewGrpcMetaData(md metadata.MD) *GrpcMetaData {
	return &GrpcMetaData{md: md}
}

func (g GrpcMetaData) Get(key string) string {
	arr := g.md[key]
	if len(arr) == 0 {
		return ""
	}
	return arr[0]
}

func (g GrpcMetaData) Set(key string, value string) {
	g.md.Set(key, value)
}

func (g GrpcMetaData) Keys() []string {
	keys := make([]string, 0, len(g.md))
	for k := range g.md {
		keys = append(keys, k)
	}
	return keys
}

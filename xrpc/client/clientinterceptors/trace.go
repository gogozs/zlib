package clientinterceptors

import (
	"context"

	"github.com/gogozs/zlib/xtrace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TraceClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
	var span trace.Span
	ctx, span = xtrace.StartSpan(ctx, method)
	defer func() { xtrace.End(span, err) }()

	md := xtrace.InjectGrpcTraceInfo(ctx)
	return invoker(metadata.NewOutgoingContext(ctx, md), method, req, reply, cc, opts...)
}

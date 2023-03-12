package serverinterceptors

import (
	"context"

	"github.com/gogozs/zlib/xtrace"
	"google.golang.org/grpc"
)

func TraceServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
	resp interface{}, err error) {
	ctx = xtrace.ExtractGrpcTraceInfo(ctx)
	ctx, span := xtrace.StartSpan(ctx, info.FullMethod)
	defer xtrace.End(span, err)

	return handler(ctx, req)
}

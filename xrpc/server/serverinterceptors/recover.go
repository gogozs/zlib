package serverinterceptors

import (
	"context"

	"github.com/gogozs/zlib/tools"
	"github.com/gogozs/zlib/xerr"
	"google.golang.org/grpc"
)

// UnaryRecoverInterceptor catches panics in processing unary requests and recovers.
func UnaryRecoverInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
	resp interface{}, err error) {
	defer tools.HandleCrash(func(r error) { err = xerr.ToGrpcError(r) })

	return handler(ctx, req)
}

// StreamRecoverInterceptor catches panics in processing stream requests and recovers.
func StreamRecoverInterceptor(svr any, stream grpc.ServerStream, _ *grpc.StreamServerInfo,
	handler grpc.StreamHandler) (err error) {
	defer tools.HandleCrash(func(r error) { err = xerr.ToGrpcError(r) })

	return handler(svr, stream)
}

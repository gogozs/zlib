package serverinterceptors

import (
	"context"
	"github.com/gogozs/zlib/tools"
	"github.com/gogozs/zlib/xerr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryErrorInterceptor ...
func UnaryErrorInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
	resp interface{}, err error) {
	defer tools.HandleCrash(func(r error) { err = xerr.ToGrpcError(r) })

	resp, err = handler(ctx, req)

	return resp, transError(err)
}

// StreamErrorInterceptor ...
func StreamErrorInterceptor(svr any, stream grpc.ServerStream, _ *grpc.StreamServerInfo,
	handler grpc.StreamHandler) (err error) {
	err = handler(svr, stream)
	return transError(err)
}

func transError(err error) error {
	xe, ok := err.(*xerr.XError)
	if !ok {
		return err
	}
	return status.Error(codes.Code(xe.Code()), xe.Error())
}

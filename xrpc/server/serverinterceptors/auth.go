package serverinterceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// refer to https://dev.to/techschoolguru/use-grpc-interceptor-for-authorization-with-jwt-1c5h

type (
	AuthInterceptor struct {
		tokenValidator TokenValidator
	}

	TokenValidator interface {
		Verify(token string) error
	}
)

func (a *AuthInterceptor) Unary(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if err = a.auth(ctx, info.FullMethod); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (a *AuthInterceptor) Stream(svr any, stream grpc.ServerStream, _ *grpc.StreamServerInfo,
	handler grpc.StreamHandler) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := a.auth(ss.Context(), info.FullMethod); err != nil {
			return err
		}

		return handler(ss, stream)
	}
}

func (a *AuthInterceptor) auth(ctx context.Context, method string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values, ok := md["Authorization"]
	if !ok || len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "token is not provided")
	}

	token := values[0]
	if err := a.tokenValidator.Verify(token); err != nil {
		return status.Errorf(codes.Unauthenticated, "token is not provided")
	}

	return nil
}

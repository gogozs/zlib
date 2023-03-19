package serverinterceptors

import (
	"context"

	"github.com/gogozs/zlib/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type (
	AuthInterceptor struct {
		authValidator AuthValidator
	}

	AuthValidator interface {
		Verify(ctx context.Context, token string) (UserInfo, error)
	}

	UserInfo interface {
		UserID() uint64
	}
)

func NewAuthInterceptor(authValidator AuthValidator) *AuthInterceptor {
	return &AuthInterceptor{
		authValidator: authValidator,
	}
}

func (a *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if err = a.auth(ctx, info.FullMethod); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (a *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := a.auth(ss.Context(), info.FullMethod); err != nil {
			return err
		}

		return handler(srv, ss)
	}
}

func (a *AuthInterceptor) auth(ctx context.Context, method string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values, ok := md["authorization"]
	if !ok || len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "token is not provided")
	}

	token := values[0]
	user, err := a.authValidator.Verify(ctx, token)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "token is not provided")
	}
	auth.WithAuth(ctx, user)
	return nil
}

func ParseUserDetails(ctx context.Context) UserInfo {
	value := auth.ParseAuth(ctx)
	return value.(UserInfo)
}

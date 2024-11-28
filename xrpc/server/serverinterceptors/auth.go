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

// wrappedServerStream is a wrapper around grpc.ServerStream that allows changing the context.
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context returns the new context.
func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

// WrapServerStream returns a new grpc.ServerStream with the provided context.
func WrapServerStream(ctx context.Context, ss grpc.ServerStream) grpc.ServerStream {
	return &wrappedServerStream{
		ServerStream: ss,
		ctx:          ctx,
	}
}

func NewAuthInterceptor(authValidator AuthValidator) *AuthInterceptor {
	return &AuthInterceptor{
		authValidator: authValidator,
	}
}

func (a *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx, err = a.auth(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (a *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx, err := a.auth(ss.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		// Wrap the original ServerStream with the new context
		wrappedSS := WrapServerStream(ctx, ss)
		return handler(srv, wrappedSS)
	}
}

func (a *AuthInterceptor) auth(ctx context.Context, method string) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values, ok := md["authorization"]
	if !ok || len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "token is not provided")
	}

	token := values[0]
	user, err := a.authValidator.Verify(ctx, token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "token is invalid")
	}
	return auth.WithAuth(ctx, user), nil
}

func ParseUserDetails(ctx context.Context) UserInfo {
	value := auth.ParseAuth(ctx)
	return value.(UserInfo)
}

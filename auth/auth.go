package auth

import "context"

type authKey struct {
}

func WithAuth(ctx context.Context, auth interface{}) context.Context {
	return context.WithValue(ctx, authKey{}, auth)
}

func ParseAuth(ctx context.Context) interface{} {
	return ctx.Value(authKey{})
}

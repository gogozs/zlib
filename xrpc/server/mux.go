package server

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

func NewMuxServer() *runtime.ServeMux {
	return runtime.NewServeMux(
		runtime.WithRoutingErrorHandler(handleRoutingError),
		runtime.WithErrorHandler(errorHandler),
	)
}

func handleRoutingError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, httpStatus int) {
	if httpStatus != http.StatusMethodNotAllowed {
		runtime.DefaultRoutingErrorHandler(ctx, mux, marshaler, w, r, httpStatus)
		return
	}

	// Use HTTPStatusError to customize the DefaultHTTPErrorHandler status code
	err := &runtime.HTTPStatusError{
		HTTPStatus: httpStatus,
		Err:        status.Error(codes.Unimplemented, http.StatusText(httpStatus)),
	}

	runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
}

func errorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, writer http.ResponseWriter, request *http.Request, err error) {
	// check if it is a custom error
	if code := status.Code(err); code > codes.Unauthenticated || code < codes.OK {
		err = &runtime.HTTPStatusError{
			HTTPStatus: http.StatusBadRequest,
			Err:        err,
		}
	}
	runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, writer, request, err)
}

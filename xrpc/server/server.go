package server

import (
	"context"
	"net"
	"net/http"

	"github.com/gogozs/zlib/tools"

	"github.com/gogozs/zlib/xrpc/server/serverinterceptors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type (
	RpcServer struct {
		options *serverOptions
	}
	RegisterFn     func(server *grpc.Server)
	RegisterHttpFn func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

	serverOptions struct {
		addr                     string
		registerFn               RegisterFn
		httpAddr                 string
		httpRegisterFn           RegisterHttpFn
		grpcOptions              []grpc.ServerOption
		unaryServerInterceptors  []grpc.UnaryServerInterceptor
		streamServerInterceptors []grpc.StreamServerInterceptor
	}

	Option func(o *serverOptions)
)

func NewRpcServer(addr string, registerFn RegisterFn, opts ...Option) *RpcServer {
	options := &serverOptions{
		addr:                     addr,
		registerFn:               registerFn,
		unaryServerInterceptors:  buildDefaultUnaryInterceptors(),
		streamServerInterceptors: buildDefaultStreamInterceptors(),
	}
	for _, o := range opts {
		o(options)
	}
	return &RpcServer{options: options}
}

func (s *RpcServer) Start() error {
	var errChan chan error
	if s.options.httpAddr != "" {
		tools.SafeGo(func() {
			errChan <- s.startHttpServer()
		})
	}
	tools.SafeGo(func() {
		errChan <- s.startServer()
	})

	return <-errChan
}

func (s *RpcServer) startServer() error {
	lis, err := net.Listen("tcp", s.options.addr)
	if err != nil {
		return err
	}

	unaryInterceptorOption := grpc.ChainUnaryInterceptor(s.options.unaryServerInterceptors...)
	streamInterceptorOption := grpc.ChainStreamInterceptor(s.options.streamServerInterceptors...)

	s.options.grpcOptions = append(s.options.grpcOptions, unaryInterceptorOption, streamInterceptorOption)
	server := grpc.NewServer(s.options.grpcOptions...)
	s.options.registerFn(server)
	return server.Serve(lis)
}

func (s *RpcServer) startHttpServer() error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := s.options.httpRegisterFn(context.Background(), mux, s.options.addr, opts); err != nil {
		return err
	}
	return http.ListenAndServe(s.options.httpAddr, mux)
}

func WithGrpcOption(options ...grpc.ServerOption) Option {
	return func(o *serverOptions) {
		o.grpcOptions = append(o.grpcOptions, options...)
	}
}

func WithUnaryServerOption(interceptors ...grpc.UnaryServerInterceptor) Option {
	return func(o *serverOptions) {
		o.unaryServerInterceptors = append(o.unaryServerInterceptors, interceptors...)
	}
}

func WithStreamServerOption(interceptors ...grpc.StreamServerInterceptor) Option {
	return func(o *serverOptions) {
		o.streamServerInterceptors = append(o.streamServerInterceptors, interceptors...)
	}
}

func WithHttpServer(httpAddr string, fn RegisterHttpFn) Option {
	return func(o *serverOptions) {
		o.httpAddr = httpAddr
		o.httpRegisterFn = fn
	}
}

func buildDefaultUnaryInterceptors() []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		serverinterceptors.UnaryRecoverInterceptor,
	}
}

func buildDefaultStreamInterceptors() []grpc.StreamServerInterceptor {
	return []grpc.StreamServerInterceptor{
		serverinterceptors.StreamRecoverInterceptor,
	}
}

func WithTLSOptions(certFile, keyFile string) grpc.ServerOption {
	creds, _ := credentials.NewServerTLSFromFile(certFile, keyFile)
	return grpc.Creds(creds)
}

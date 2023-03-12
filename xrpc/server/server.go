package server

import (
	"net"

	"github.com/gogozs/zlib/xrpc/server/serverinterceptors"
	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
)

type (
	RpcServer struct {
		addr       string
		registerFn RegisterFn
		options    *serverOptions
	}
	RegisterFn func(server *grpc.Server)

	serverOptions struct {
		grpcOptions              []grpc.ServerOption
		unaryServerInterceptors  []grpc.UnaryServerInterceptor
		streamServerInterceptors []grpc.StreamServerInterceptor
	}

	Option func(o *serverOptions)
)

func NewRpcServer(addr string, registerFn RegisterFn, opts ...Option) *RpcServer {
	options := &serverOptions{
		unaryServerInterceptors:  buildDefaultUnaryInterceptors(),
		streamServerInterceptors: buildDefaultStreamInterceptors(),
	}
	for _, o := range opts {
		o(options)
	}
	return &RpcServer{addr: addr, registerFn: registerFn, options: options}
}

func (s *RpcServer) Start() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	unaryInterceptorOption := grpc.ChainUnaryInterceptor(s.options.unaryServerInterceptors...)
	streamInterceptorOption := grpc.ChainStreamInterceptor(s.options.streamServerInterceptors...)

	s.options.grpcOptions = append(s.options.grpcOptions, unaryInterceptorOption, streamInterceptorOption)
	server := grpc.NewServer(s.options.grpcOptions...)
	s.registerFn(server)

	return server.Serve(lis)
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

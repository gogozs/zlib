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
		options    []grpc.ServerOption
	}
	RegisterFn func(server *grpc.Server)
)

func NewRpcServer(addr string, registerFn RegisterFn, options ...grpc.ServerOption) *RpcServer {
	return &RpcServer{addr: addr, registerFn: registerFn, options: options}
}

func (s *RpcServer) Start() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	unaryInterceptorOption := grpc.ChainUnaryInterceptor(s.buildUnaryInterceptors()...)
	streamInterceptorOption := grpc.ChainStreamInterceptor(s.buildStreamInterceptors()...)

	s.options = append(s.options, unaryInterceptorOption, streamInterceptorOption)

	server := grpc.NewServer(s.options...)
	s.registerFn(server)

	return server.Serve(lis)
}

func (s *RpcServer) buildUnaryInterceptors() []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		serverinterceptors.UnaryRecoverInterceptor,
	}
}

func (s *RpcServer) buildStreamInterceptors() []grpc.StreamServerInterceptor {
	return []grpc.StreamServerInterceptor{
		serverinterceptors.StreamRecoverInterceptor,
	}
}

func WithTLSOptions(certFile, keyFile string) grpc.ServerOption {
	creds, _ := credentials.NewServerTLSFromFile(certFile, keyFile)
	return grpc.Creds(creds)
}

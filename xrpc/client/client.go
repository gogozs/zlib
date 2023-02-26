package client

import (
	"github.com/gogozs/zlib/xerr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type (
	RpcClient struct {
		conn *grpc.ClientConn
	}
)

func NewRpcClient(addr string, options ...grpc.DialOption) (*RpcClient, error) {
	c := &RpcClient{}
	if err := c.dail(addr, options...); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *RpcClient) Conn() *grpc.ClientConn {
	return c.conn
}

func (c *RpcClient) dail(addr string, options ...grpc.DialOption) error {
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, options...)
	if err != nil {
		return xerr.Errorf("did not connect: %v", err)
	}
	c.conn = conn
	return nil
}

func WithTLSOptions(certFile string) grpc.DialOption {
	creds, _ := credentials.NewClientTLSFromFile(certFile, "")
	return grpc.WithTransportCredentials(creds)
}

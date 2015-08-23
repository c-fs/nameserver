package client

import (
	pb "github.com/c-fs/nameserver/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Client struct {
	grpcConn   *grpc.ClientConn
	nameClient pb.NameClient
}

func New(address string) (*Client, error) {
	conn, err := grpc.Dial(address)
	if err != nil {
		return nil, err
	}

	return &Client{grpcConn: conn, nameClient: pb.NewNameClient(conn)}, nil
}

func (c *Client) FetchDisks(ctx context.Context) ([]*pb.DiskInfo, error) {
	reply, err := c.nameClient.FetchDisks(ctx, &pb.FetchDisksRequest{})

	if err != nil {
		return nil, err
	}
	return reply.Disks, nil
}

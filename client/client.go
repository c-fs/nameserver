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

func (c *Client) FetchDisks(ctx context.Context) (map[string]string, error) {
	reply, err := c.nameClient.FetchDisks(ctx, &pb.FetchDisksRequest{})

	if err != nil {
		return nil, err
	}

	diskMap := make(map[string]string)
	for _, disk := range reply.Disks {
		diskMap[disk.Name] = disk.Remote
	}
	return diskMap, nil
}

package main

import (
	"github.com/BurntSushi/toml"
	pb "github.com/c-fs/nameserver/proto"
	"github.com/c-fs/nameserver/server/config"
	"github.com/qiniu/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
)

type server struct {
	registeredDisks []*pb.DiskInfo
}

func NewServer() *server {
	return &server{registeredDisks: make([]*pb.DiskInfo, 0, 26)}
}

func (s *server) FetchDisks(context.Context, *pb.FetchDisksRequest) (*pb.FetchDisksReply, error) {
	return &pb.FetchDisksReply{Disks: s.registeredDisks}, nil
}

func main() {
	configfn := "nameserver.conf"
	data, err := ioutil.ReadFile(configfn)
	if err != nil {
		log.Fatalf("server: cannot load configuration file[%s] (%v)", configfn, err)
	}

	var conf config.Server
	if _, err := toml.Decode(string(data), &conf); err != nil {
		log.Fatalf("server: configuration file[%s] is not valid (%v)", configfn, err)
	}
	server := NewServer()
	for i, v := range conf.Disks {
		log.Infof("Adding %v to disks", v)
		server.registeredDisks = append(server.registeredDisks, &conf.Disks[i])
	}
	log.Infof("server: starting server...")

	lis, err := net.Listen("tcp", net.JoinHostPort(conf.Bind, conf.Port))

	if err != nil {
		log.Fatalf("server: failed to listen: %v", err)
	}

	log.Infof("server: listening on %s", net.JoinHostPort(conf.Bind, conf.Port))

	s := grpc.NewServer()
	pb.RegisterNameServer(s, server)
	log.Infof("server: ready to serve clients")
	s.Serve(lis)
}

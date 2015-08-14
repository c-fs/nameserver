package config

import (
	pb "github.com/c-fs/nameserver/proto"
)

type Server struct {
	Disks []pb.DiskInfo
	Port  string
	Bind  string
}

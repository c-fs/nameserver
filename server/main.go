package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/c-fs/cfs/client"
	pb "github.com/c-fs/nameserver/proto"
	"github.com/c-fs/nameserver/server/config"
	"github.com/qiniu/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
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
	// TODO: support server list
	apiServer := flag.String("kubernetes-api-server", "", "Kubernetes api server address. If set, it uses disks registered in kubernetes instead of disks in config file.")
	flag.Parse()

	configfn := "nameserver.conf"
	data, err := ioutil.ReadFile(configfn)
	if err != nil {
		log.Fatalf("server: cannot load configuration file[%s] (%v)", configfn, err)
	}

	var conf config.Server
	if _, err := toml.Decode(string(data), &conf); err != nil {
		log.Fatalf("server: configuration file[%s] is not valid (%v)", configfn, err)
	}

	var disks []pb.DiskInfo
	if *apiServer != "" {
		var err error
		// TODO: customized labelSelector
		disks, err = getDisksFromKubernetesAPI(*apiServer, "cfs")
		if err != nil {
			log.Fatalf("server: failed to get disks from --kubernetes-api-server (%v)", err)
		}
	} else {
		disks = conf.Disks
	}

	server := NewServer()
	for i, v := range disks {
		log.Infof("Adding %v to disks", v)
		server.registeredDisks = append(server.registeredDisks, &disks[i])
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

func getDisksFromKubernetesAPI(serverURL, label string) ([]pb.DiskInfo, error) {
	resp, err := getPodsResponse(serverURL, label)
	if err != nil {
		return nil, err
	}
	log.Printf("%s", resp)
	ips, err := parsePodIPs(resp)
	if err != nil {
		return nil, err
	}
	log.Printf("%v", ips)
	return getDisks(ips)
}

func getPodsResponse(serverURL, label string) ([]byte, error) {
	resp, err := http.Get(serverURL + "/api/v1/pods?labelSelector=name=" + label)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

type podsResponse struct {
	Items []struct {
		Status struct {
			PodIP string `json:"podIP"`
		} `json:"status"`
	} `json:"items"`
}

func parsePodIPs(resp []byte) ([]string, error) {
	var pods podsResponse
	if err := json.Unmarshal(resp, &pods); err != nil {
		return nil, err
	}
	var ips []string
	for _, item := range pods.Items {
		ips = append(ips, item.Status.PodIP)
	}
	return ips, nil
}

func getDisks(ips []string) ([]pb.DiskInfo, error) {
	var disks []pb.DiskInfo
	for _, ip := range ips {
		// TODO: customized cfs port
		c, err := client.New(0x5678, ip+":15524")
		if err != nil {
			return nil, err
		}
		ds, err := c.Disks(context.TODO())
		for _, d := range ds {
			disks = append(disks, pb.DiskInfo{
				Name:   d.Name,
				Remote: ip,
				Port:   15524,
			})
		}
	}
	return disks, nil
}

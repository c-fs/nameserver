package handlers

import (
	"net/http"

	"github.com/c-fs/nameserver/client"
	pb "github.com/c-fs/nameserver/proto"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

var (
	NameServerAddr = ""
	nameClient     *client.Client
)

func InitHandlers(addr string) error {
	var err error
	nameClient, err = client.New(addr)
	if err != nil {
		return err
	}
	return nil
}

type FetchDisksRequest struct {
	pb.FetchDisksRequest
}

func FetchDisks(c *gin.Context) {
	var req FetchDisksRequest
	c.Bind(&req)

	reply, err := nameClient.FetchDisks(context.TODO())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pb.FetchDisksReply{Disks: reply})
}

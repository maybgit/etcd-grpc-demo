package services

import (
	"context"
	"etcd-grpc-service/proto"
	"fmt"
)

type HelloEtcd struct{}

func (c *HelloEtcd) SayHello(ctx context.Context, in *proto.NameRequest) (*proto.BaseResponse, error) {
	fmt.Println(in.Name)
	return &proto.BaseResponse{Code: 200, Message: "Hello " + in.Name}, nil
}

package services

import (
	"etcd-grpc-service/proto"
	"context"
)

type HelloEtcd struct{ }

func (c *HelloEtcd) SayHello(ctx context.Context, in *proto.NameRequest) (*proto.BaseResponse, error) {
	return &proto.BaseResponse{Code: 200,Message: "Hello "+in.Name},nil
}
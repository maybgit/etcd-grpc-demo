package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/coreos/etcd/clientv3"
	etcdnaming "go.etcd.io/etcd/clientv3/naming"

	"etcd-grpc-client/proto"

	"google.golang.org/grpc"
)

func main() {
	serviceName := "etcd-grpc-service"                        //服务名称
	cli, err := clientv3.NewFromURL("http://10.1.1.248:2379") //连接etcd服务
	if err != nil {
		fmt.Println(err)
	}
	i := 1
	//从服务注册列表中去发现服务并创建连接
	for {
		func() {
			r := &etcdnaming.GRPCResolver{Client: cli}
			b := grpc.RoundRobin(r)
			if conn, gerr := grpc.Dial(serviceName, grpc.WithBalancer(b), grpc.WithBlock(), grpc.WithInsecure()); gerr != nil {
				fmt.Println(gerr)
			} else {
				defer conn.Close()
				client := proto.NewHelloEtcdClient(conn)
				if res, err := client.SayHello(context.TODO(), &proto.NameRequest{Name: strconv.Itoa(i)}); err != nil {
					fmt.Println(err)
				} else {
					println(res.Code, res.Message)
				}
			}
			b.Close()
		}()
		i++
		time.Sleep(time.Millisecond * 100)
		if i == 100 {
			break
		}
	}
}

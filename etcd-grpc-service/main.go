package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"etcd-grpc-service/proto"
	"etcd-grpc-service/services"

	"github.com/coreos/etcd/clientv3"
	etcdnaming "go.etcd.io/etcd/clientv3/naming"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/naming"
)

func main() {
	serviceName := "etcd-grpc-service" //服务名称
	serviceAddr := "10.1.1.54:100"     //服务地址

	cli, err := clientv3.NewFromURL("http://10.1.1.248:2379") //连接etcd服务
	if err != nil {
		fmt.Println(err)
	}
	r := &etcdnaming.GRPCResolver{Client: cli}
	r.Update(context.TODO(), serviceName, naming.Update{Op: naming.Add, Addr: serviceAddr}) //服务注册

	//启动grpc服务
	mux := http.NewServeMux()
	server := grpc.NewServer()
	proto.RegisterHelloEtcdServer(server, &services.HelloEtcd{})

	//启动grpc服务
	if err := http.ListenAndServe(serviceAddr, func() http.Handler {
		return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
				server.ServeHTTP(w, r)
			} else {
				mux.ServeHTTP(w, r)
			}
		}), &http2.Server{})

	}()); err != nil {
		fmt.Println(err)
	}
}

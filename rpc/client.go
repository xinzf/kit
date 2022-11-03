package rpc

import (
	"context"
	etcd_client "github.com/rpcxio/rpcx-etcd/client"
	"github.com/smallnest/rpcx/client"
	rpcx_protocol "github.com/smallnest/rpcx/protocol"
	"github.com/xinzf/kit/container/kcfg"
	"github.com/xinzf/kit/container/kvar"
	"github.com/xinzf/kit/klog"
)

func Call(service, method string, req, rsp any) error {

	basePath := kcfg.Get[string]("base")
	if basePath == "" {
		basePath = "gokit"
	}

	registry := kcfg.Get[string]("registry")
	if registry == "" {
		registry = "etcd"
	}

	etcdAddrs := kvar.New(kcfg.Get[[]any]("rpc.etcd.addrs")).Strings()
	if etcdAddrs == nil || len(etcdAddrs) == 0 {
		klog.Fatal("Missing etcd addrs")
	}

	option := client.DefaultOption
	option.SerializeType = rpcx_protocol.JSON
	d, _ := etcd_client.NewEtcdV3Discovery(basePath, service, etcdAddrs, false, nil)
	xclient := client.NewXClient(service, client.Failover, client.RoundRobin, d, option)
	defer xclient.Close()

	return xclient.Call(context.Background(), method, req, rsp)
}

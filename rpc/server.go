package rpc

import (
	"context"
	"fmt"
	"github.com/liushuochen/gotable"
	"github.com/liushuochen/gotable/cell"
	"github.com/liushuochen/gotable/table"
	"github.com/rcrowley/go-metrics"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"github.com/xinzf/kit/container/kcfg"
	"github.com/xinzf/kit/container/kvar"
	"github.com/xinzf/kit/klog"
	"time"
)

var _handlers []*handler

func init() {
	_handlers = []*handler{}
}

func Register(handlers ...interface{}) {
	for _, _hdl := range handlers {
		if hdl := newHandler(_hdl); hdl != nil {
			_handlers = append(_handlers, hdl)
		}
	}
}

func Run(ctx context.Context, before ...func(serv *server.Server) error) {
	serviceName := kcfg.Get[string]("rpc.name")
	if serviceName == "" {
		klog.Fatal("Missing rpc service's port")
	}

	port := kcfg.Get[int]("rpc.port")
	if port == 0 {
		klog.Fatal("Missing rpc service's port")
	}
	addr := fmt.Sprintf("localhost:%d", port)

	var (
		tb  *table.Table
		err error
	)

	tb, err = gotable.Create("#", "Service", "Method", "Handler")
	if err != nil {
		klog.Panicf("Create print table failed: %s", err.Error())
	}

	serv := server.NewServer()
	if err = register(serv, addr); err != nil {
		klog.Fatal(err.Error())
	}

	if len(before) > 0 {
		for _, f := range before {
			if err := f(serv); err != nil {
				klog.Fatal(err.Error())
			}
		}
	}

	tb.Align("#", cell.AlignLeft)
	tb.Align("Service", cell.AlignLeft)
	tb.Align("Method", cell.AlignLeft)
	tb.Align("Handler", cell.AlignLeft)
	tb.SetColumnColor("#", gotable.Underline, gotable.Write, gotable.NoneBackground)
	tb.SetColumnColor("Service", gotable.Underline, gotable.Write, gotable.NoneBackground)
	tb.SetColumnColor("Method", gotable.Underline, gotable.Write, gotable.NoneBackground)
	tb.SetColumnColor("Handler", gotable.Underline, gotable.Write, gotable.NoneBackground)

	num := 0
	for _, h := range _handlers {
		for _, method := range h.methods {
			num++
			_ = tb.AddRow([]string{
				fmt.Sprintf("%d", num),
				fmt.Sprintf("%s.%s", serviceName, h.serviceName),
				method,
				fmt.Sprintf("%s/%s", h.pkgPath, h.handlerName),
			})
		}
		err = serv.RegisterName(fmt.Sprintf("%s.%s", serviceName, h.serviceName), h.hdl, "")
		if err != nil {
			klog.Fatal(err.Error())
		}
	}

	fmt.Println()
	fmt.Printf("[SERVER] server listen on port: %d, Total: %d\n", port, num)
	fmt.Println(tb)
	err = serv.Serve("tcp", addr)
	if err != nil {
		klog.Fatal(err.Error())
	}
}

func register(serv *server.Server, addr string) error {

	basePath := kcfg.Get[string]("base")
	if basePath == "" {
		basePath = "gokit"
	}

	registry := kcfg.Get[string]("registry")
	if registry == "" {
		registry = "etcd"
	}

	switch registry {
	case "etcd":
		etcdAddrs := kvar.New(kcfg.Get[[]any]("rpc.etcd.addrs")).Strings()
		if etcdAddrs == nil || len(etcdAddrs) == 0 {
			klog.Fatal("Missing etcd addrs")
		}
		reg := &serverplugin.EtcdV3RegisterPlugin{
			ServiceAddress: "tcp@" + addr,
			EtcdServers:    etcdAddrs,
			BasePath:       basePath,
			Metrics:        metrics.NewRegistry(),
			UpdateInterval: time.Second * 30,
		}
		if err := reg.Start(); err != nil {
			return err
		}

		serv.Plugins.Add(reg)
	}
	return nil
}

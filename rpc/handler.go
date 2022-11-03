package rpc

import (
	"github.com/xinzf/kit/klog"
	kitServer "github.com/xinzf/kit/server"
	"reflect"
)

func newHandler(hand any) *handler {
	refValue := reflect.ValueOf(hand)
	refType := reflect.TypeOf(hand)
	handlerName := refType.Elem().Name()
	serviceName := handlerName
	if refType.Implements(reflect.TypeOf(new(kitServer.HandlerName)).Elem()) {
		aliasMethod := refValue.MethodByName("HandlerName")
		alisValue := aliasMethod.Call([]reflect.Value{})
		serviceName = alisValue[0].Interface().(string)
	}

	for _, h := range _handlers {
		if h.serviceName == serviceName {
			klog.Warnf("There are multiple handlers with the same service name: %s", serviceName)
			return nil
		}
	}

	pkgPath := refType.Elem().PkgPath()
	methods := make([]string, 0)

	for i := 0; i < refValue.NumMethod(); i++ {
		methodName := refType.Method(i).Name
		method := refValue.Method(i)

		if method.Type().NumIn() != 3 || method.Type().NumOut() != 1 {
			continue
		}

		var (
			c      = method.Type().In(0)
			req    = method.Type().In(1)
			rsp    = method.Type().In(2)
			output = method.Type().Out(0)
		)

		if c.String() != "context.Context" {
			klog.Warnf("the context argument of %s.%s is not context.Context", handlerName, methodName)
			continue
		}

		if req.Kind() != reflect.Ptr && req.Kind() != reflect.Interface {
			klog.Warnf("the request argument of %s.%s is not a pointer", handlerName, methodName)
			continue
		}

		if rsp.Kind() != reflect.Interface && rsp.Kind() != reflect.Pointer {
			klog.Warnf("the response argument of %s.%s is not a pointer", handlerName, methodName)
			continue
		}

		if output.String() != "error" {
			klog.Warnf("the output error of %s.%s is not error type", handlerName, methodName)
		}

		methods = append(methods, methodName)
	}

	if len(methods) == 0 {
		return nil
	}

	return &handler{
		hdl:         hand,
		pkgPath:     pkgPath,
		handlerName: handlerName,
		serviceName: serviceName,
		methods:     methods,
	}
}

type handler struct {
	hdl         any
	pkgPath     string
	handlerName string
	serviceName string
	methods     []string
}

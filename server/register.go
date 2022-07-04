package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"strings"
)

type HandlerGroup struct {
	parent      *HandlerGroup
	path        string
	middlewares []gin.HandlerFunc
	handlers    []*handler
	subGroups   []*HandlerGroup
}

var groups []*HandlerGroup

func init() {
	groups = make([]*HandlerGroup, 0)
}

func (this *HandlerGroup) getPath() string {
	currPath := strings.TrimLeft(strings.TrimRight(this.path, "/"), "/")
	if this.parent != nil {
		parentPath := strings.TrimLeft(strings.TrimRight(this.parent.getPath(), "/"), "/")
		return fmt.Sprintf("%s/%s", parentPath, currPath)
	}
	return currPath
}

func (this *HandlerGroup) Group(path string, middlewares ...gin.HandlerFunc) *HandlerGroup {
	for _, subGroup := range this.subGroups {
		if subGroup.path == path {
			return subGroup
		}
	}

	g := &HandlerGroup{
		path:        path,
		middlewares: middlewares,
		handlers:    []*handler{},
		subGroups:   []*HandlerGroup{},
		parent:      this,
	}
	this.subGroups = append(this.subGroups, g)
	return g
}

func (this *HandlerGroup) Register(apiHandler ...interface{}) *HandlerGroup {
	for _, handler := range apiHandler {
		refValue := reflect.ValueOf(handler)
		refType := reflect.TypeOf(handler)
		handlerName := refType.Elem().Name()

		alisName := ""
		if refType.Implements(reflect.TypeOf(new(HandlerName)).Elem()) {
			aliasMethod := refValue.MethodByName("HandlerName")
			alisValue := aliasMethod.Call([]reflect.Value{})
			alisName = alisValue[0].Interface().(string)
		}

		for i := 0; i < refValue.NumMethod(); i++ {
			methodName := refType.Method(i).Name

			for _, h := range this.handlers {
				if h.handlerName == handlerName && h.methodName == methodName {
					continue
				}
			}

			hdl := newHandler(refType.Elem().PkgPath(), handlerName, methodName, refValue.Method(i), alisName)
			if hdl == nil {
				continue
			}

			this.handlers = append(this.handlers, hdl)
		}
	}
	return this
}

func Group(path string, middlewares ...gin.HandlerFunc) *HandlerGroup {
	g := &HandlerGroup{
		path:        path,
		middlewares: middlewares,
		handlers:    []*handler{},
		subGroups:   []*HandlerGroup{},
	}
	groups = append(groups, g)
	return g
}

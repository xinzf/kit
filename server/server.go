package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/os/glog"
	"github.com/liushuochen/gotable"
	"github.com/liushuochen/gotable/cell"
	"github.com/liushuochen/gotable/table"
	"github.com/xinzf/kit/container/kcfg"
	"strings"
)

func Run(ctx context.Context, before ...func() error) {
	if len(before) > 0 {
		for _, f := range before {
			if err := f(); err != nil {
				glog.Panic(err)
			}
		}
	}

	debug := kcfg.Get[bool]("server.debug")
	port := kcfg.Get[int]("server.port")
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	g := gin.Default()

	var (
		tb  *table.Table
		err error
	)

	tb, err = gotable.Create("#", "URI", "Handler", "Method")
	if err != nil {
		glog.Panicf("Create print table failed: %s", err.Error())
	}

	tb.Align("#", cell.AlignLeft)
	tb.Align("URI", cell.AlignLeft)
	tb.Align("Handler", cell.AlignLeft)
	tb.Align("Method", cell.AlignLeft)
	tb.SetColumnColor("#", gotable.Underline, gotable.Write, gotable.NoneBackground)
	tb.SetColumnColor("URI", gotable.Underline, gotable.Write, gotable.NoneBackground)
	tb.SetColumnColor("Handler", gotable.Underline, gotable.Write, gotable.NoneBackground)
	tb.SetColumnColor("Method", gotable.Underline, gotable.Write, gotable.NoneBackground)

	num := 0
	var listen = func(ginGroup *gin.RouterGroup, _group *HandlerGroup) {}
	listen = func(ginGroup *gin.RouterGroup, _group *HandlerGroup) {
		for _, h := range _group.handlers {
			path := ""
			if _group.path == "" {
				_group.path = "/"
			}
			if _group.path == "/" {
				path = strings.TrimLeft(h.getBindPath(), "/")
			} else {
				path = h.getBindPath()
			}

			evalPath := func() string {
				str := fmt.Sprintf("%s%s", _group.getPath(), path)
				if !strings.HasPrefix(str, "/") {
					str = "/" + str
				}
				return str
			}

			num++
			_ = tb.AddRow([]string{
				fmt.Sprintf("%d", num),
				evalPath(),
				fmt.Sprintf("%s/%s", h.pkgPath, h.handlerName),
				h.methodName,
			})

			ginGroup.OPTIONS(path, func(c *gin.Context) {
				c.JSON(200, nil)
			})
			ginGroup.POST(path, h.handlerFunc)
		}
		for _, subGroup := range _group.subGroups {
			_subGinGroup := ginGroup.Group(subGroup.path, subGroup.middlewares...)
			listen(_subGinGroup, subGroup)
		}
	}

	for _, _group := range groups {
		listen(g.Group(_group.path, _group.middlewares...), _group)
	}

	fmt.Println()
	fmt.Printf("[SERVER] server listen on port: %d, Total: %d\n", port, num)
	fmt.Println(tb)
	_ = g.Run(fmt.Sprintf(":%d", port))
}

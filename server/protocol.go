package server

import (
	"github.com/gin-gonic/gin"
)

//
//type RequestVerify interface {
//	Verify() (int, error)
//}
//
//type RequestCacheSetter interface {
//	CacheId() (string, time.Duration)
//}
//
//type RequestCacheDeleter interface {
//	DeleteCachePatterns() (patterns []string)
//}

type AfterRequestInterface interface {
	AfterRequest(c *gin.Context)
}

type AfterResponseInterface interface {
	AfterResponse(c *gin.Context, rsp ResponseInterface)
}

type ResponseInterface interface {
	GetData() any
	GetStatus() int
	GetMessage() string
}

type Response struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   any    `json:"data"`
}

func (r Response) GetData() any {
	return r.Data
}

func (r Response) GetStatus() int {
	return r.Status
}

func (r Response) GetMessage() string {
	return r.Msg
}

type HandlerName interface {
	HandlerName() string
}

type HandlerPath interface {
	Paths() map[string][]string
}

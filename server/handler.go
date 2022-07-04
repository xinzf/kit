package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xinzf/kit/klog"
	"reflect"
	"time"
)

func newHandler(pkgPath, handlerName, methodName string, fun reflect.Value, aliasName string) (hdl *handler) {
	if fun.Type().NumIn() > 3 || fun.Type().NumIn() < 2 || fun.Type().NumOut() < 1 || fun.Type().NumOut() > 2 {
		return
	}

	var (
		c          reflect.Type
		req        reflect.Type
		rsp        reflect.Type
		outputCode reflect.Type
		outputErr  reflect.Type
		isUpload   bool
	)

	if fun.Type().NumIn() == 2 {
		req = fun.Type().In(0)
		rsp = fun.Type().In(1)
	} else {
		c = fun.Type().In(0)
		req = fun.Type().In(1)
		rsp = fun.Type().In(2)
	}

	if fun.Type().NumOut() == 2 {
		outputCode = fun.Type().Out(0)
		outputErr = fun.Type().Out(1)
	} else {
		outputErr = fun.Type().Out(0)
	}

	if req.Kind() != reflect.Interface && req.Kind() != reflect.Ptr {
		klog.Warnf("the request argument of %s.%s is not a pointer", handlerName, methodName)
		return
	}

	if rsp.Kind() != reflect.Interface && rsp.Kind() != reflect.Map {
		if rsp.Kind() != reflect.Ptr {
			klog.Warnf("the response argument of %s.%s is not a pointer", handlerName, methodName)
			return
		}
	}

	if fun.Type().NumIn() == 3 && c.String() != "*gin.Context" {
		klog.Warnf("the context argument of %s.%s is not *gin.Context", handlerName, methodName)
		return
	}

	if fun.Type().NumOut() == 2 && outputCode.String() != "int" {
		klog.Warnf("the output code of %s.%s is not a int", handlerName, methodName)
		return
	}

	if outputErr.String() != "error" {
		klog.Warnf("the output error of %s.%s is not error type", handlerName, methodName)
		return
	}

	if req.String() == "*server.UploadRequest" {
		isUpload = true
	}

	hdl = &handler{
		pkgPath:     pkgPath,
		uri:         "",
		handlerName: handlerName,
		aliasName:   aliasName,
		methodName:  methodName,
		fun:         fun,
		c:           c,
		req:         req,
		rsp:         rsp,
		inNum:       fun.Type().NumIn(),
		outNum:      fun.Type().NumOut(),
		outputCode:  outputCode,
		outputErr:   outputErr,
		isUpload:    isUpload,
	}

	return
}

type handler struct {
	pkgPath     string
	uri         string
	handlerName string
	aliasName   string
	methodName  string
	fun         reflect.Value
	c           reflect.Type
	req         reflect.Type
	rsp         reflect.Type
	inNum       int
	outNum      int
	outputCode  reflect.Type
	outputErr   reflect.Type
	isUpload    bool
	cache       struct {
		key      string
		duration time.Duration
	}
}

func (this *handler) getBindPath() string {
	var path string
	if this.aliasName != "" {
		path = fmt.Sprintf("/%s/%s", snakeString(this.aliasName), snakeString(this.methodName))
	} else {
		path = fmt.Sprintf("/%s/%s", snakeString(this.handlerName), snakeString(this.methodName))
	}
	return path
}

func (this *handler) handlerFunc(c *gin.Context) {
	var response Response
	defer func() {
		if response.Data == nil {
			response.Data = map[string]interface{}{}
		}
	}()

	var req reflect.Value
	if this.req.Kind() == reflect.Interface {
		req = reflect.New(this.req)
	} else {
		req = reflect.New(this.req.Elem())
	}

	var uploadFile *UploadRequest
	if this.isUpload == false {
		contentType := c.ContentType()
		if contentType == "text/xml" {
			if err := c.BindXML(req.Interface()); err != nil {
				response.Status = 400
				response.Msg = err.Error()
				c.AbortWithStatusJSON(200, response)
				return
			}
		} else if contentType == "application/json" {
			if err := c.BindJSON(req.Interface()); err != nil {
				response.Status = 400
				response.Msg = err.Error()
				c.AbortWithStatusJSON(200, response)
				return
			}
		}
	} else {
		uploadFile = &UploadRequest{ctx: c}
	}

	//this.requestVerify(req, c)
	this.afterRequest(req, c)
	if c.IsAborted() {
		return
	}

	//this.requestCache(req, c)
	//if c.IsAborted() {
	//	return
	//}

	var rsp reflect.Value
	if this.rsp.String() == "interface {}" {
		rsp = reflect.ValueOf(map[string]any{})
	} else {
		rsp = reflect.New(this.rsp.Elem())
	}

	values := make([]reflect.Value, 0)
	if this.inNum == 2 {
		if this.isUpload {
			values = this.fun.Call([]reflect.Value{
				reflect.ValueOf(uploadFile),
				rsp,
			})
		} else {
			values = this.fun.Call([]reflect.Value{
				req, rsp,
			})
		}
	} else {
		if this.isUpload {
			values = this.fun.Call([]reflect.Value{
				reflect.ValueOf(c),
				reflect.ValueOf(uploadFile),
				rsp,
			})
		} else {
			values = this.fun.Call([]reflect.Value{
				reflect.ValueOf(c),
				req, rsp,
			})
		}
	}

	var (
		code int = 0
		err  error
	)
	if this.outNum == 2 {
		code, _ = values[0].Interface().(int)
		err, _ = values[1].Interface().(error)
	} else {
		err, _ = values[0].Interface().(error)
	}

	if err != nil && code == 0 {
		code = 500
	}

	if code != 0 || err != nil {
		response.Status = code
		response.Msg = err.Error()
		c.AbortWithStatusJSON(200, response)
		return
	}

	response.Data = rsp.Interface()
	output := Response{
		Status: 0,
		Msg:    "",
		Data:   response.Data,
	}

	c.JSON(200, output)
	this.afterResponse(req, c, output)
	//
	//if this.cache.key != "" && response.Data != nil {
	//	ref := reflect.ValueOf(response.Data)
	//	switch ref.Kind() {
	//	case reflect.Array, reflect.Map, reflect.Slice, reflect.Pointer, reflect.Struct:
	//		bts, err := jsoniter.Marshal(response.Data)
	//		if err == nil {
	//			err = cache.Redis().Set(context.TODO(), this.cache.key, bts, this.cache.duration).Err()
	//		}
	//	default:
	//		err = cache.Redis().Set(context.TODO(), this.cache.key, response.Data, this.cache.duration).Err()
	//	}
	//	if err != nil {
	//		klog.Errof("put request cache failed: %s", err.Error())
	//	} else {
	//		klog.Debugf("save request cache: %s success", this.cache.key)
	//	}
	//}
	//
	//if err = this.deleteCache(req); err != nil {
	//	klog.Errof("delete request cache failed: %s", err.Error())
	//}

	return
}

func (this *handler) afterRequest(reqVal reflect.Value, c *gin.Context) {
	if reqVal.Type().Implements(reflect.TypeOf(new(AfterRequestInterface)).Elem()) == false {
		return
	}
	//if this.implementAfterRequest == false {
	//    return
	//}
	//if _, found := this.req.MethodByName("Verify"); !found {
	//	return
	//}
	fn := reqVal.MethodByName("AfterRequest")
	//initFunc := reqVal.MethodByName("Verify")
	//if initFunc.Type().NumIn() != 0 || initFunc.Type().NumOut() != 2 {
	//	return
	//}
	//if initFunc.Type().Out(0).Kind() != reflect.Int && initFunc.Type().Out(1).String() != "error" {
	//	return
	//}

	fn.Call([]reflect.Value{reflect.ValueOf(c)})
	//values := fn.Call([]reflect.Value{reflect.ValueOf(c)})
	//status := values[0].Interface().(int)
	//msg := ""
	//if status > 0 && !values[1].IsNil() {
	//	msg = values[1].Interface().(error).Error()
	//}
	//if values[0].Interface().(int) != 0 {
	//	c.AbortWithStatusJSON(200, Response{
	//		Status: status,
	//		Msg:    msg,
	//		Data:   map[string]interface{}{},
	//	})
	//}
}

func (this *handler) afterResponse(reqVal reflect.Value, c *gin.Context, rsp ResponseInterface) {
	if reqVal.Type().Implements(reflect.TypeOf(new(AfterResponseInterface)).Elem()) == false {
		return
	}

	fn := reqVal.MethodByName("AfterResponse")
	fn.Call([]reflect.Value{reflect.ValueOf(c), reflect.ValueOf(rsp)})
}

//
//func (this *handler) requestCache(reqVal reflect.Value, c *gin.Context) {
//	if this.isUpload {
//		return
//	}
//
//	if _, found := this.req.MethodByName("CacheId"); !found {
//		return
//	}
//	initFunc := reqVal.MethodByName("CacheId")
//	if initFunc.Type().NumOut() != 2 {
//		return
//	}
//	if initFunc.Type().Out(0).Kind() != reflect.String || initFunc.Type().Out(1).Kind() != reflect.Int64 {
//		return
//	}
//
//	values := initFunc.Call([]reflect.Value{})
//	this.cache.key = values[0].Interface().(string)
//	this.cache.duration = values[1].Interface().(time.Duration)
//	if this.cache.key == "" {
//		return
//	}
//
//	hit := false
//	bts, err := cache.Redis().Get(context.TODO(), this.cache.key).Bytes()
//	if err != nil {
//		if err.Error() == redis.Nil.Error() {
//			err = nil
//		}
//	} else {
//		hit = true
//	}
//
//	if err != nil {
//		c.AbortWithStatusJSON(200, Response{
//			Status: 500,
//			Msg:    fmt.Sprintf("get request cache failed: %s", err.Error()),
//			Data:   map[string]interface{}{},
//		})
//		return
//	}
//
//	if !hit {
//		return
//	}
//
//	var data any
//	err = jsoniter.Unmarshal(bts, &data)
//	if err != nil {
//		klog.Errof("get request cache failed: %s", err.Error())
//		return
//	}
//	c.AbortWithStatusJSON(200, Response{
//		Status: 0,
//		Msg:    "",
//		Data:   data,
//	})
//	klog.Debugf("hit request cache: %s", this.cache.key)
//}
//
//func (this *handler) deleteCache(reqVal reflect.Value) error {
//	if this.isUpload {
//		return nil
//	}
//
//	if _, found := this.req.MethodByName("DeleteCachePatterns"); !found {
//		return nil
//	}
//	initFunc := reqVal.MethodByName("DeleteCachePatterns")
//	if initFunc.Type().NumOut() != 1 {
//		return nil
//	}
//	if initFunc.Type().Out(0).Kind() != reflect.Slice {
//		return nil
//	}
//
//	values := initFunc.Call([]reflect.Value{})
//	if values[0].Interface() == nil {
//		return nil
//	}
//
//	patterns := kvar.New(values[0].Interface()).Strings()
//	if len(patterns) == 0 {
//		return nil
//	}
//
//	keys := make([]string, 0)
//	for _, s := range patterns {
//		if strings.Contains(s, "*") {
//			keys = append(keys, cache.Redis().Keys(context.TODO(), s).Val()...)
//		} else {
//			keys = append(keys, s)
//		}
//	}
//
//	if len(keys) == 0 {
//		return nil
//	}
//	return cache.Redis().Del(context.TODO(), keys...).Err()
//}

package utils

import "reflect"

func Clone(src interface{}) interface{} {
	var val interface{}
	typ := reflect.TypeOf(src)
	v := reflect.ValueOf(src)
	if v.Kind() == reflect.Ptr {
		typ = typ.Elem()               //获取源实际类型(否则为指针类型)
		dst := reflect.New(typ).Elem() //创建对象
		data := reflect.ValueOf(src)   //源数据值
		data = data.Elem()             //源数据实际值（否则为指针）
		dst.Set(data)                  //设置数据
		dst = dst.Addr()               //创建对象的地址（否则返回值）
		val = dst.Interface()
	} else {
		dst := reflect.New(typ).Elem() //创建对象
		data := reflect.ValueOf(src)   //源数据值
		dst.Set(data)                  //设置数据
		val = dst.Interface()
	}
	return val
}

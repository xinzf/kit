package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/xinzf/kit/container/kcfg"
	"github.com/xinzf/kit/klog"
)

var inited = false
var _redisClient *redis.Client
var _ctx context.Context = context.TODO()

func connect() error {
	if inited {
		return nil
	}
	inited = true
	addr := kcfg.Get[string]("cache.host")
	pswd := kcfg.Get[string]("cache.pswd")
	db := kcfg.Get[int]("cache.db")

	_redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pswd,
		DB:       db,
	})
	if err := _redisClient.Ping(_ctx).Err(); err != nil {
		return fmt.Errorf("connect redis server failed: %s", err.Error())
	}
	klog.Args("host", addr, "db", db).Info("Redis connect success!")
	return nil
}

func Redis() *redis.Client {
	if _redisClient == nil || _redisClient.Ping(context.TODO()).Err() != nil {
		if err := connect(); err != nil {
			panic(any(err))
		}
	}
	//if !inited {
	//	if err := connect(); err != nil {
	//		panic(any(err))
	//	}
	//}

	if _redisClient == nil {
		klog.Panic("Redis has not connect, please connect redis first")
	}
	return _redisClient
}

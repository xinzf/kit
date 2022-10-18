package kcfg

import (
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"github.com/xinzf/kit/container/klist"
	"github.com/xinzf/kit/container/kmap"
)

type config struct {
	v *viper.Viper
}

var inited bool
var cfg *config

func init() {
	if inited {
		return
	}
	inited = true
	cfg = &config{v: viper.New()}
	cfg.v.SetDefault("name", "gokit")
	cfg.v.SetDefault("server.port", 8080)
	cfg.v.SetDefault("server.debug", true)
	cfg.v.SetDefault("logger.level", "debug")
	cfg.v.SetDefault("logger.type", "text")
	cfg.v.SetDefault("logger.stack", "panic")
}

func New(cfgFile string) error {
	cfg.v.SetConfigType("yaml")
	cfg.v.SetConfigFile(cfgFile)
	if err := cfg.v.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func Set[T any](key string, value T) {
	cfg.v.Set(key, value)
}

func Get[T any](key string) (result T) {
	if val := cfg.v.Get(key); val == nil {
		return
	} else {
		return val.(T)
	}
}

func GetMap[T any](key string) (result *kmap.Map[string, T]) {
	if val := cfg.v.GetStringMap(key); val == nil {
		return
	} else {
		data := make(map[string]T)
		for k, v := range val {
			data[k] = v.(T)
		}
		mp := kmap.New[string, T]()
		mp.MergeMap(data)
		return mp
	}
}

func GetList[T any](key string) (result *klist.List[T]) {
	result = klist.New[T]()
	val := cfg.v.Get(key)
	if val == nil {
		return
	}
	if vals, err := cast.ToSliceE(val); err != nil {
		return
	} else {
		for _, v := range vals {
			result.Add(v.(T))
		}
		return
	}
}

package kmap

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/emirpasic/gods/maps/hashmap"
	jsoniter "github.com/json-iterator/go"
	"github.com/xinzf/kit/container/kvar"
	"golang.org/x/exp/constraints"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"strings"
)

type Map[K constraints.Complex | constraints.Ordered, V any] struct {
	mp *hashmap.Map
}

func New[K constraints.Complex | constraints.Ordered, V any](mp ...map[K]V) *Map[K, V] {
	values := make(map[K]V)
	if len(mp) > 0 {
		values = mp[0]
	}

	_mp := hashmap.New()
	for k, v := range values {
		_mp.Put(k, v)
	}
	return &Map[K, V]{mp: _mp}
}

func (this *Map[K, V]) Set(key K, value V) {
	this.mp.Put(key, value)
}

func (this *Map[K, V]) Get(key K) (value V, found bool) {
	var val interface{}
	val, found = this.mp.Get(key)
	if !found || val == nil {
		return
	}
	value = val.(V)
	return
}

func (this *Map[K, V]) GetOrSet(key K, value V) (val V) {
	var found bool
	val, found = this.Get(key)
	if found {
		return
	}

	val = value
	this.Set(key, val)
	return
}

func (this *Map[K, V]) GetOrSetFunc(key K, fn func() V) (val V) {
	var found bool
	val, found = this.Get(key)
	if found {
		return
	}

	val = fn()
	this.Set(key, val)
	return
}

func (this *Map[K, V]) Remove(keys ...K) {
	for _, key := range keys {
		this.mp.Remove(key)
	}
}

func (this *Map[K, V]) IsEmpty() bool {
	if this.Size() == 0 {
		return true
	}

	empty := true
	for _, val := range this.mp.Values() {
		if val != nil {
			empty = false
			break
		}
	}
	return empty
}

func (this *Map[K, V]) Size() int {
	return this.mp.Size()
}

func (this *Map[K, V]) Keys() []K {
	keys := make([]K, 0, this.Size())
	for _, v := range this.mp.Keys() {
		keys = append(keys, v.(K))
	}
	return keys
}

func (this *Map[K, V]) Values() (values []V) {
	values = make([]V, 0, this.Size())
	for _, v := range this.mp.Values() {
		if v == nil {
			var alias V
			_ = kvar.New(nil).Convert(&alias)
			values = append(values, alias)
		} else {
			values = append(values, v.(V))
		}
	}
	return
}

func (this *Map[K, V]) FilterEmpty() {
	keys := this.mp.Keys()
	for _, key := range keys {
		value, _ := this.mp.Get(key)
		if value == nil {
			this.mp.Remove(key)
		}
	}
}

func (this *Map[K, V]) Clone() *Map[K, V] {
	newMp := make(map[K]V)
	bts, err := this.mp.ToJSON()
	if err == nil {
		_ = jsoniter.Unmarshal(bts, &newMp)
		return New[K, V](newMp)
	}
	return New[K, V]()
	//jsoniter.Unmarshal()
	//newMp := New[K, V]()
	//keys := this.mp.Keys()
	//for _, key := range keys {
	//	newKey := utils.Clone(key)
	//	value, _ := this.mp.Get(key)
	//	newValue := utils.Clone(value)
	//	newMp.mp.Put(newKey, newValue)
	//}
	//return newMp
}

func (this *Map[K, V]) Clear() {
	this.mp.Clear()
}

func (this *Map[K, V]) Find(fn func(item V) bool) (K, V, bool) {
	keys := this.mp.Keys()

	for _, key := range keys {
		value, _ := this.mp.Get(key)
		var alias V
		if value == nil {
			value = kvar.New(nil).Convert(&alias)
		} else {
			alias = value.(V)
		}

		if fn(alias) {
			return key.(K), alias, true
		}
	}

	var (
		k K
		v V
	)
	_ = kvar.New(nil).Convert(&k)
	_ = kvar.New(nil).Convert(&v)
	return k, v, false
}

func (this *Map[K, V]) Each(fn func(key K, value V)) {
	keys := this.mp.Keys()
	for _, key := range keys {
		value, _ := this.mp.Get(key)
		fn(key.(K), value.(V))
	}
}

func (this *Map[K, V]) Iterator(fn func(key K, value V) error) error {
	keys := this.mp.Keys()
	for _, key := range keys {
		value, _ := this.mp.Get(key)
		if err := fn(key.(K), value.(V)); err != nil {
			return err
		}
	}
	return nil
}

func (this *Map[K, V]) Map() map[K]V {
	mp := make(map[K]V)
	for _, key := range this.mp.Keys() {
		value, _ := this.mp.Get(key)
		mp[key.(K)] = value.(V)
	}
	return mp
}

func (this *Map[K, V]) Merge(mp *Map[K, V]) {
	mp.Each(func(key K, value V) {
		this.Set(key, value)
	})
}

func (this *Map[K, V]) MergeMap(mp map[K]V) {
	for k, v := range mp {
		this.Set(k, v)
	}
}

// Value return json value, implement driver.Valuer interface
func (m Map[K, V]) Value() (driver.Value, error) {
	bytes, err := m.mp.ToJSON()
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}

// Scan value into Jsonb, implements sql.Scanner interface
func (m *Map[K, V]) Scan(value interface{}) error {
	*m = Map[K, V]{hashmap.New()}
	if value == nil {
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	return m.mp.FromJSON(bytes)
}

// MarshalJSON to output non base64 encoded []byte
func (m Map[K, V]) MarshalJSON() ([]byte, error) {
	return m.mp.ToJSON()
}

// UnmarshalJSON to deserialize []byte
func (m *Map[K, V]) UnmarshalJSON(b []byte) error {
	return m.mp.FromJSON(b)
}

func (this *Map[K, V]) String() string {
	data, _ := this.mp.ToJSON()
	return string(data)
}

// GormDataType gorm common data type
func (Map[K, V]) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (Map[K, V]) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	case "sqlserver":
		return "NVARCHAR(MAX)"
	}
	return ""
}

func (km Map[K, V]) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if km.mp.Size() == 0 {
		return gorm.Expr("{}")
	}

	data, _ := km.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}

	return gorm.Expr("?", string(data))
}

package klist

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/emirpasic/gods/lists/arraylist"
	jsoniter "github.com/json-iterator/go"
	"github.com/xinzf/kit/container/kvar"
	"github.com/xinzf/kit/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"strings"
)

type ElementEqual interface {
	Equal(val any) bool
}

type List[T any] struct {
	list *arraylist.List
}

func (this *List[T]) Add(values ...T) {
	if len(values) == 0 {
		return
	}
	equal := func(val T) bool {
		//value := values[0]
		//refType := reflect.TypeOf(value)
		//if !refType.Implements(reflect.TypeOf(new(ElementEqual)).Elem()) {
		//	return false
		//}
		//
		//method := reflect.ValueOf(val).MethodByName("Equal")
		//idx, _ := this.list.Find(func(_ int, elem interface{}) bool {
		//	_values := method.Call([]reflect.Value{reflect.ValueOf(elem)})
		//	return _values[0].Interface().(bool)
		//})
		//if idx == -1 {
		//	return false
		//}
		return false
	}

	for _, value := range values {
		if equal(value) {
			continue
		} else {
			this.list.Add(value)
		}
	}
}

func (this *List[T]) Get(index int) (value T, found bool) {
	var val interface{}
	val, found = this.list.Get(index)
	if found {
		if val == nil {
			var alias T
			_ = kvar.New(val).Convert(&alias)
			return alias, true
		} else {
			value = val.(T)
		}
	}
	return
}

func (this *List[T]) Remove(index int) {
	this.list.Remove(index)
}

func (this *List[T]) RemoveValue(value T) {
	this.list.Remove(this.list.IndexOf(value))
}

func (this *List[T]) RemoveByFilter(filter func(value T) bool) {
	values := this.list.Values()
	for idx, val := range values {
		value := val.(T)
		if filter(value) {
			this.list.Remove(idx)
		}
	}
}

func (this *List[T]) Contains(values ...T) bool {
	vals := make([]interface{}, 0)
	{
		for _, value := range values {
			vals = append(vals, value)
		}
	}
	return this.list.Contains(vals...)
}

func (this *List[T]) List(filter ...func(elem T) bool) []T {
	values := this.list.Values()
	elements := make([]T, 0)
	for _, value := range values {
		if len(filter) > 0 {
			if filter[0](value.(T)) {
				elements = append(elements, value.(T))
			}
		} else {
			elements = append(elements, value.(T))
		}
	}
	return elements
}

func (this *List[T]) Each(fn func(idx int, elem T)) {
	this.list.Each(func(index int, value interface{}) {
		elem := value.(T)
		fn(index, elem)
	})
}

func (this *List[T]) Iterator(fn func(idx int, elem T) error) error {
	values := this.list.Values()
	for i, value := range values {
		elem := value.(T)
		if err := fn(i, elem); err != nil {
			return err
		}
	}
	return nil
}

func (this *List[T]) Find(fn func(ele T) bool) (idx int, value T) {
	values := this.list.Values()
	var val interface{}
	for idx, val = range values {
		if val == nil {
			var alias T
			_ = kvar.New(val).Convert(&alias)
			value = alias
		} else {
			value = val.(T)
		}
		if fn(value) {
			return
		}
	}

	idx = -1
	return
}

func (this *List[T]) Select(fn ...func(idx int, ele T) bool) *List[T] {
	newList := New[T]()
	values := this.list.Values()
	for i, value := range values {
		var val T
		if value == nil {
			var alias T
			_ = kvar.New(value).Convert(&alias)
			val = alias
		} else {
			val = value.(T)
		}
		if len(fn) > 0 {
			if fn[0](i, val) {
				newList.Add(val)
			}
		} else {
			newList.Add(val)
		}
	}
	return newList
}

func (this *List[T]) Clone(fn ...func(idx int, ele T) bool) *List[T] {
	newList := New[T]()
	values := this.list.Values()
	for i, value := range values {
		elem := utils.Clone(value).(T)
		if len(fn) > 0 {
			if fn[0](i, elem) {
				newList.Add(elem)
			}
		} else {
			newList.Add(elem)
		}
	}
	return newList
}

func (this *List[T]) Size(filter ...func(idx int, elem T) bool) int {
	return this.Select(filter...).list.Size()
}

func (this *List[T]) IsEmpty() bool {
	return this.list.Empty()
}

func (this *List[T]) Clear() {
	this.list.Clear()
}

func (this *List[T]) IndexOf(value T) int {
	return this.list.IndexOf(value)
}

// Sort if a < b return true, or return false
func (this *List[T]) Sort(comparator func(a, b T) (asc bool)) {
	this.list.Sort(func(a, b interface{}) int {
		if comparator(a.(T), b.(T)) {
			return -1
		}
		return 1
	})
}

// Swap ???????????????????????????
func (this *List[T]) Swap(i, j int) {
	this.list.Swap(i, j)
}

func (this *List[T]) Insert(index int, values ...T) {
	vals := make([]interface{}, 0)
	{
		for _, value := range values {
			vals = append(vals, value)
		}
	}
	this.list.Insert(index, vals...)
}

func (this *List[T]) Reverse() {
	values := this.list.Values()
	this.list.Clear()
	for i := len(values) - 1; i >= 0; i-- {
		this.list.Add(values[i].(T))
	}
}

func (this *List[T]) MoveAfter(target, curr T) {
	idx := this.IndexOf(target)
	if idx == -1 {
		return
	}

	currIdx := this.IndexOf(curr)
	if currIdx == -1 {
		return
	}

	this.Remove(currIdx)
	this.Insert(idx+1, curr)
}

func (this *List[T]) MoveBefore(target, curr T) {
	idx := this.IndexOf(target)
	if idx == -1 {
		return
	}

	currIdx := this.IndexOf(curr)
	if currIdx == -1 {
		return
	}

	this.Remove(currIdx)
	this.Insert(idx, curr)
}

func (this *List[T]) MoveToBack(value T) {
	if idx := this.IndexOf(value); idx != -1 {
		this.Remove(idx)
		this.list.Insert(this.list.Size(), value)
	}
}

func (this *List[T]) MoveToFront(value T) {
	if idx := this.IndexOf(value); idx != -1 {
		this.Remove(idx)
		this.list.Insert(0, value)
	}
}

func (this *List[T]) PushBack(value T) {
	this.list.Insert(this.list.Size(), value)
}

func (this *List[T]) PushFront(value T) {
	this.list.Insert(0, value)
}

func (this *List[T]) PushBackList(list *List[T]) {
	list.Each(func(idx int, elem T) {
		this.list.Add(elem)
	})
}

func (this *List[T]) PushFrontList(list *List[T]) {
	frontValues := list.list.Values()
	values := this.list.Values()
	finalValues := append(frontValues, values...)
	this.list.Clear()
	this.list.Add(finalValues...)
}

func (this *List[T]) First() T {
	elem, _ := this.Get(0)
	return elem
}

func (this *List[T]) Last() T {
	elem, _ := this.Get(this.Size() - 1)
	return elem
}

func (this *List[T]) Set(index int, value T) {
	if index < 0 {
		this.list.Add(value)
	} else {
		this.list.Set(index, value)
	}
}

// Intersect ???????????????set?????????others??????????????????????????????others ?????????????????????????????????
func (this *List[T]) Intersect(others ...*List[T]) (newList *List[T]) {
	newList = New[T]()
	this.Each(func(_ int, a T) {
		same := true
		for _, other := range others {
			idx, _ := other.Find(func(b T) bool {
				return kvar.New(a).Equal(kvar.New(b))
			})
			if idx == -1 {
				same = false
				break
			}
		}
		if same {
			newList.Add(a)
		}
	})
	return
}

// Diff ???????????????set????????????others??????????????????????????????others?????????????????????????????????
func (this *List[T]) Diff(others ...*List[T]) (newList *List[T]) {
	newList = New[T]()

	this.Each(func(_ int, a T) {
		same := true
		for _, other := range others {
			if idx, _ := other.Find(func(b T) bool {
				return kvar.New(a).Equal(kvar.New(b))
			}); idx != -1 {
				same = false
				break
			}
		}
		if !same {
			newList.Add(a)
		}
	})

	return
}

// Union ???????????????set?????????others??????????????????????????????
func (this *List[T]) Union(others ...*List[T]) (newList *List[T]) {
	newList = New[T]()
	others = append(others, this)
	for _, other := range others {
		other.Each(func(_ int, a T) {
			idx, _ := newList.Find(func(b T) bool {
				return kvar.New(a).Equal(kvar.New(b))
			})
			if idx == -1 {
				newList.Add(a)
			}
		})
	}
	return
}

//Distinct ??????
func (this *List[T]) Distinct() (newList *List[T]) {
	newList = New[T]()
	this.Each(func(_ int, a T) {
		idx, _ := newList.Find(func(b T) bool {
			return kvar.New(a).Equal(kvar.New(b))
		})
		if idx == -1 {
			newList.Add(a)
		}
	})
	return newList
}

func (s List[T]) MarshalJSON() ([]byte, error) {
	return s.list.ToJSON()
}

func (s *List[T]) UnmarshalJSON(bytes []byte) error {
	values := make([]T, 0)
	if err := jsoniter.Unmarshal(bytes, &values); err != nil {
		return err
	}

	list := arraylist.New()
	for _, value := range values {
		list.Add(value)
	}
	*s = List[T]{list: list}
	//*s = List[T]{list: arraylist.New()}
	return s.list.FromJSON(bytes)
}

func (s *List[T]) UnmarshalBinary(bytes []byte) error {
	values := make([]T, 0)
	if err := jsoniter.Unmarshal(bytes, &values); err != nil {
		return err
	}

	list := arraylist.New()
	for _, value := range values {
		list.Add(value)
	}
	*s = List[T]{list: list}
	//*s = List[T]{list: arraylist.New()}
	return s.list.FromJSON(bytes)
}

func (s *List[T]) Scan(value interface{}) error {
	*s = List[T]{list: arraylist.New()}
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
	return s.list.FromJSON(bytes)
}

func (s List[T]) Value() (driver.Value, error) {
	bytes, err := s.list.ToJSON()
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}

func (s List[T]) String() string {
	return fmt.Sprintf("list%+v", s.list.Values())
}

// GormDataType gorm common data type
func (List[T]) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (List[T]) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (s List[T]) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if s.list.Size() == 0 {
		switch db.Dialector.Name() {
		case "mysql":
			return gorm.Expr("[]")
		default:
			return gorm.Expr("NULL")
		}
	}

	data, _ := s.list.ToJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}

	return gorm.Expr("?", string(data))
}

func New[T any](values ...T) *List[T] {
	l := &List[T]{list: arraylist.New()}
	l.Add(values...)
	return l
}

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
	"reflect"
	"strings"
)

type List[T any] struct {
	list *arraylist.List
}

func New[T any](values ...T) *List[T] {
	l := &List[T]{list: arraylist.New()}
	l.Add(values...)
	return l
}

func (this *List[T]) Add(values ...T) {
	if len(values) == 0 {
		return
	}
	for _, value := range values {
		this.list.Add(value)
	}
}

func (this *List[T]) Get(index int) (value T, found bool) {
	var val interface{}
	val, found = this.list.Get(index)
	if found {
		if val != nil {
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
		var alias T
		if val != nil {
			alias = val.(T)
		}
		if filter(alias) {
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
	if filter == nil || len(filter) == 0 {
		filter = []func(elem T) bool{func(elem T) bool {
			return true
		}}
	}

	values := this.list.Values()
	elements := make([]T, 0)
	for _, value := range values {
		var alias T
		if value != nil {
			alias = value.(T)
		}
		if filter[0](alias) {
			elements = append(elements, alias)
		}
	}
	return elements
}

func (this *List[T]) Each(fn func(idx int, elem T)) {
	this.list.Each(func(index int, value interface{}) {
		var alias T
		if value != nil {
			alias = value.(T)
		}
		fn(index, alias)
	})
}

func (this *List[T]) Iterator(fn func(idx int, elem T) error) error {
	values := this.list.Values()
	for i, value := range values {
		var alias T
		if value != nil {
			alias = value.(T)
		}

		if err := fn(i, alias); err != nil {
			return err
		}
	}
	return nil
}

func (this *List[T]) Find(fn func(ele T) bool) (int, T) {
	values := this.list.Values()
	for idx, val := range values {
		var value T
		if val != nil {
			value = val.(T)
		}
		if fn(value) {
			return idx, value
		}
	}

	var value T
	return -1, value
}

func (this *List[T]) Select(filter ...func(idx int, ele T) bool) *List[T] {
	if filter == nil || len(filter) == 0 {
		filter = []func(idx int, ele T) bool{func(idx int, ele T) bool {
			return true
		}}
	}

	newList := New[T]()
	values := this.list.Values()
	for i, value := range values {
		var alias T
		if value != nil {
			alias = value.(T)
		}
		if filter[0](i, alias) {
			newList.Add(alias)
		}
	}
	return newList
}

func (this *List[T]) Clone(filter ...func(idx int, ele T) bool) *List[T] {
	if filter == nil || len(filter) == 0 {
		filter = []func(idx int, ele T) bool{func(idx int, ele T) bool {
			return true
		}}
	}

	newList := New[T]()
	values := this.list.Values()
	for i, value := range values {
		var elem T
		if value != nil {
			elem = utils.Clone(value).(T)
		}

		if filter[0](i, elem) {
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

// Swap 交换两个元素的位置
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

func (this *List[T]) Pop() T {
	elem := this.First()
	this.Remove(0)
	return elem
}

func (this *List[T]) Set(index int, value T) {
	if index < 0 {
		this.list.Add(value)
	} else {
		this.list.Set(index, value)
	}
}

//Distinct 去重
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
			return gorm.Expr("CAST('[]' AS JSON)")
		case "postgres":
			return gorm.Expr("CAST('[]' AS JSONB)")
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
	case "postgres":
		return gorm.Expr("CAST(? AS JSONB)", string(data))
	}

	return gorm.Expr("?", string(data))
}

func (this *List[T]) Convert(target any) error {
	if reflect.Indirect(reflect.ValueOf(target)).Kind() != reflect.Pointer {
		return errors.New("convert failed, target must be a pointer")
	}

	btes, err := this.list.ToJSON()
	if err != nil {
		return err
	}

	return jsoniter.Unmarshal(btes, target)
}

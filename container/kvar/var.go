package kvar

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/golang-module/carbon/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"reflect"
	"strings"
	"time"
)

type Var struct {
	value        any
	kind         reflect.Kind
	isComparable bool
}

func New(val any) *Var {
	v := &Var{
		value: val,
		kind:  reflect.Indirect(reflect.ValueOf(val)).Kind(),
	}

	switch v.kind {
	case reflect.String,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Float32,
		reflect.Float64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Bool:
		v.isComparable = true
	default:
	}

	return v
}

func (this *Var) IsNil() bool {
	return this.value == nil
}

func (this *Var) IsInt() bool {
	switch this.kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	default:
		return false
	}
}

func (this *Var) IsFloat() bool {
	switch this.kind {
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func (this *Var) IsString() bool {
	return this.kind == reflect.String
}

func (this *Var) IsSlice() bool {
	return this.kind == reflect.Slice
}

func (this *Var) IsMap() bool {
	return this.kind == reflect.Map
}

func (this *Var) IsStruct() bool {
	return this.kind == reflect.Struct
}

func (this *Var) IsPointer() bool {
	return this.kind == reflect.Pointer
}

func (this *Var) IsUint() bool {
	switch this.kind {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}

func (this *Var) IsEmpty() bool {
	if this.IsNil() {
		return true
	} else if this.IsString() {
		return this.String() == ""
	} else if this.IsInt() {
		return this.Int64() == 0
	} else if this.IsSlice() {
		return len(this.Interfaces()) == 0
	} else if this.Bool() {
		return this.Bool() == false
	} else if this.IsFloat() {
		return this.Float64() == 0
	} else if this.IsMap() {
		return len(this.Map()) == 0
	} else if this.IsStruct() {
		tpe := reflect.TypeOf(this.value)
		return tpe.NumField() == 0 && tpe.NumMethod() == 0
	} else if this.IsPointer() {
		return reflect.ValueOf(this.value).IsNil()
	} else if this.IsUint() {
		return this.Uint64() == 0
	} else {
		return true
	}
}

func (this *Var) String() string {
	if this.isComparable {
		return cast.ToString(this.value)
	}
	return ""
}

func (this *Var) Strings() []string {
	if !this.IsSlice() {
		return []string{}
	}

	return cast.ToStringSlice(this.value)
}

func (this *Var) Int() int {
	if this.isComparable {
		return cast.ToInt(this.value)
	}
	return 0
}

func (this *Var) Ints() []int {
	if this.IsSlice() == false {
		return []int{}
	}

	return cast.ToIntSlice(this.value)
}

func (this *Var) int8() int8 {
	if this.isComparable {
		return cast.ToInt8(this.value)
	}
	return 0
}

func (this *Var) Int8s() []int8 {
	if this.IsSlice() == false {
		return []int8{}
	}

	ints := this.Ints()
	int8s := make([]int8, 0)
	if ints != nil {
		for _, i := range ints {
			int8s = append(int8s, cast.ToInt8(i))
		}
	}
	return int8s
}

func (this *Var) Int16() int16 {
	if this.isComparable {
		return cast.ToInt16(this.value)
	}
	return 0
}

func (this *Var) Int16s() []int16 {
	if this.IsSlice() == false {
		return []int16{}
	}
	int16s := make([]int16, 0)
	ints := this.Ints()
	if ints != nil {
		for _, i := range ints {
			int16s = append(int16s, cast.ToInt16(i))
		}
	}
	return int16s
}

func (this *Var) Int32() int32 {
	if this.isComparable {
		return cast.ToInt32(this.value)
	}
	return 0
}

func (this *Var) Int32s() []int32 {
	if this.IsSlice() == false {
		return []int32{}
	}

	int32s := make([]int32, 0)
	ints := this.Ints()
	if ints != nil {
		for _, i := range ints {
			int32s = append(int32s, cast.ToInt32(i))
		}
	}
	return int32s
}

func (this *Var) Int64() int64 {
	if this.isComparable {
		return cast.ToInt64(this.value)
	}
	return 0
}

func (this *Var) Int64s() []int64 {
	if !this.IsSlice() {
		return []int64{}
	}
	int64s := make([]int64, 0)
	ints := this.Ints()
	if ints != nil {
		for _, i := range ints {
			int64s = append(int64s, cast.ToInt64(i))
		}
	}
	return int64s
}

func (this *Var) Float32() float32 {
	if this.isComparable {
		return cast.ToFloat32(this.value)
	}
	return 0
}

func (this *Var) Float32s() []float32 {
	if !this.IsSlice() {
		return []float32{}
	}
	vals := cast.ToSlice(this.value)
	if vals == nil {
		return []float32{}
	}

	floats := make([]float32, 0)
	for _, val := range vals {
		if v, err := cast.ToFloat32E(val); err == nil {
			floats = append(floats, v)
		} else {
			floats = append(floats, 0)
		}
	}
	return floats
}

func (this *Var) Float64() float64 {
	return cast.ToFloat64(this.value)
}

func (this *Var) Float64s() []float64 {
	if !this.IsSlice() {
		return []float64{}
	}

	vals := make([]any, 0)
	bts, _ := jsoniter.Marshal(this.value)
	_ = jsoniter.Unmarshal(bts, &vals)

	if vals == nil {
		return []float64{}
	}

	floats := make([]float64, 0)
	for _, val := range vals {
		floats = append(floats, New(val).Float64())
		//if v, err := cast.ToFloat64E(val); err == nil {
		//	floats = append(floats, v)
		//} else {
		//	floats = append(floats, 0)
		//}
	}
	return floats
}

func (this *Var) Uint() uint {
	return cast.ToUint(this.value)
}

func (this *Var) Uints() []uint {
	if !this.IsSlice() {
		return []uint{}
	}

	vals := cast.ToSlice(this.value)
	if vals == nil {
		return []uint{}
	}

	uints := make([]uint, 0)
	for _, val := range vals {
		if v, err := cast.ToUintE(val); err == nil {
			uints = append(uints, v)
		} else {
			uints = append(uints, 0)
		}
	}
	return uints
}

func (this *Var) Uint8() uint8 {
	return cast.ToUint8(this.value)
}

func (this *Var) Uint8s() []uint8 {
	if !this.IsSlice() {
		return []uint8{}
	}
	vals := cast.ToSlice(this.value)
	if vals == nil {
		return []uint8{}
	}

	uint8s := make([]uint8, 0)
	for _, val := range vals {
		if v, err := cast.ToUint8E(val); err == nil {
			uint8s = append(uint8s, v)
		} else {
			uint8s = append(uint8s, 0)
		}
	}
	return uint8s
}

func (this *Var) Uint16() uint16 {
	return cast.ToUint16(this.value)
}

func (this *Var) Uint16s() []uint16 {
	if !this.IsSlice() {
		return []uint16{}
	}
	vals := cast.ToSlice(this.value)
	if vals == nil {
		return []uint16{}
	}

	uint16s := make([]uint16, 0)
	for _, val := range vals {
		if v, err := cast.ToUint16E(val); err == nil {
			uint16s = append(uint16s, v)
		} else {
			uint16s = append(uint16s, 0)
		}
	}
	return uint16s
}

func (this *Var) Uint32() uint32 {
	return cast.ToUint32(this.value)
}

func (this *Var) Uint32s() []uint32 {
	if !this.IsSlice() {
		return []uint32{}
	}
	vals := cast.ToSlice(this.value)
	if vals == nil {
		return []uint32{}
	}

	uint32s := make([]uint32, 0)
	for _, val := range vals {
		if v, err := cast.ToUint32E(val); err == nil {
			uint32s = append(uint32s, v)
		} else {
			uint32s = append(uint32s, 0)
		}
	}
	return uint32s
}

func (this *Var) Uint64() uint64 {
	return cast.ToUint64(this.value)
}

func (this *Var) Uint64s() []uint64 {
	if !this.IsSlice() {
		return []uint64{}
	}
	vals := cast.ToSlice(this.value)
	if vals == nil {
		return []uint64{}
	}

	uint64s := make([]uint64, 0)
	for _, val := range vals {
		if v, err := cast.ToUint64E(val); err == nil {
			uint64s = append(uint64s, v)
		} else {
			uint64s = append(uint64s, 0)
		}
	}
	return uint64s
	//return cast.ToUint64(this.value)
}

func (this *Var) Bool() bool {
	return cast.ToBool(this.value)
}

func (this *Var) Bools() []bool {
	if !this.IsSlice() {
		return []bool{}
	}
	return cast.ToBoolSlice(this.value)
}

func (this *Var) Interface() interface{} {
	return this.value
}

func (this *Var) Interfaces() []interface{} {
	if !this.IsSlice() {
		return []interface{}{}
	}

	values := make([]any, 0)
	btes, _ := jsoniter.Marshal(this.value)
	jsoniter.Unmarshal(btes, &values)
	return values

	//return cast.ToSlice(this.value)
}

func (this *Var) Vars() []*Var {
	vals := this.Interfaces()
	if vals == nil {
		return []*Var{}
	}

	vars := make([]*Var, 0)
	for _, val := range vals {
		vars = append(vars, New(val))
	}
	return vars
}

func (this *Var) Duration() time.Duration {
	return cast.ToDuration(this.value)
}

func (this *Var) Time() time.Time {
	parseTimestamp := func() (time.Time, error) {
		timestamp := this.Int64()
		if timestamp == 0 {
			return time.Time{}, errors.New("err")
		}

		var c carbon.Carbon
		str := fmt.Sprintf("%d", timestamp)
		if len(str) == 10 {
			c = carbon.CreateFromTimestamp(timestamp, carbon.Shanghai)
		} else if len(str) == 13 {
			c = carbon.CreateFromTimestampMilli(timestamp, carbon.Shanghai)
		} else if len(str) == 16 {
			c = carbon.CreateFromTimestampMicro(timestamp, carbon.Shanghai)
		} else if len(str) == 19 {
			c = carbon.CreateFromTimestampNano(timestamp, carbon.Shanghai)
		}
		if err := c.Error; err != nil {
			return time.Time{}, err
		}
		return c.Carbon2Time(), nil
	}

	parseString := func() (time.Time, error) {
		str := this.String()
		if str == "" {
			return time.Time{}, errors.New("err")
		}
		c := carbon.Parse(str)
		if err := c.Error; err != nil {
			return time.Time{}, err
		}
		return c.Carbon2Time(), nil
	}

	if t, err := parseTimestamp(); err == nil {
		return t
	} else if t, err = parseString(); err == nil {
		return t
	} else {
		return time.Time{}
	}
}

func (this *Var) Times() []time.Time {
	vars := this.Vars()
	if vars == nil {
		return []time.Time{}
	}

	times := make([]time.Time, 0)
	for _, v := range vars {
		times = append(times, v.Time())
	}
	return times
}

func (this *Var) Bytes() []byte {
	if this.isComparable {
		str := this.String()
		return []byte(str)
	}

	bt, _ := jsoniter.Marshal(this.value)
	return bt
}

func (this *Var) Kind() reflect.Kind {
	return this.kind
}

func (this *Var) Equal(elem *Var) bool {
	if this.Kind() != elem.Kind() {
		return false
	}
	return reflect.DeepEqual(this.value, elem.value)
}

func (this *Var) Map() map[string]interface{} {
	return cast.ToStringMap(this.value)
}

func (this *Var) MapStr() map[string]string {
	return cast.ToStringMapString(this.value)
}

func (this *Var) MapInt() map[string]int {
	return cast.ToStringMapInt(this.value)
}

func (this *Var) MapInt8() map[string]int8 {
	mp := this.MapInt()
	if mp == nil {
		return map[string]int8{}
	}
	newMp := make(map[string]int8)
	{
		for s, i := range mp {
			newMp[s] = cast.ToInt8(i)
		}
	}
	return newMp
}

func (this *Var) MapInt16() map[string]int16 {
	mp := this.MapInt()
	if mp == nil {
		return map[string]int16{}
	}
	newMp := make(map[string]int16)
	{
		for s, i := range mp {
			newMp[s] = cast.ToInt16(i)
		}
	}
	return newMp
}

func (this *Var) MapInt32() map[string]int32 {
	mp := this.MapInt()
	if mp == nil {
		return map[string]int32{}
	}
	newMp := make(map[string]int32)
	{
		for s, i := range mp {
			newMp[s] = cast.ToInt32(i)
		}
	}
	return newMp
}

func (this *Var) MapInt64() map[string]int64 {
	return cast.ToStringMapInt64(this.value)
}

func (this *Var) MapFloat32() map[string]float32 {
	mp := this.Map()
	if mp == nil {
		return map[string]float32{}
	}

	newMp := make(map[string]float32)
	for s, i := range mp {
		newMp[s] = cast.ToFloat32(i)
	}
	return newMp
}

func (this *Var) MapFloat64() map[string]float64 {
	mp := this.Map()
	if mp == nil {
		return map[string]float64{}
	}

	newMp := make(map[string]float64)
	for s, i := range mp {
		newMp[s] = cast.ToFloat64(i)
	}
	return newMp
}

func (this *Var) MapBool() map[string]bool {
	return cast.ToStringMapBool(this.value)
}

func (this *Var) MapVar() map[string]*Var {
	mp := this.Map()
	newMps := make(map[string]*Var)
	for s, i := range mp {
		newMps[s] = New(i)
	}
	return newMps
}

func (this *Var) MapSlice() map[string][]interface{} {
	mps := cast.ToStringMapStringSlice(this.value)
	newMps := make(map[string][]interface{})
	for s, strings := range mps {
		newMps[s] = []interface{}{}
		{
			for _, s2 := range strings {
				newMps[s] = append(newMps[s], s2)
			}
		}
	}
	return newMps
}

func (this *Var) MapSliceStrings() map[string][]string {
	return cast.ToStringMapStringSlice(this.value)
}

func (this *Var) MapSliceVars() map[string][]*Var {
	mps := this.MapSlice()
	if mps == nil {
		return map[string][]*Var{}
	}
	vars := make(map[string][]*Var)
	for s, i := range mps {
		vars[s] = []*Var{}
		for _, i2 := range i {
			vars[s] = append(vars[s], New(i2))
		}
	}
	return vars
}

func (this *Var) Clone() *Var {
	var val interface{}
	typ := reflect.TypeOf(this.value)
	if this.kind == reflect.Ptr {
		typ = typ.Elem()                    //获取源实际类型(否则为指针类型)
		dst := reflect.New(typ).Elem()      //创建对象
		data := reflect.ValueOf(this.value) //源数据值
		data = data.Elem()                  //源数据实际值（否则为指针）
		dst.Set(data)                       //设置数据
		dst = dst.Addr()                    //创建对象的地址（否则返回值）
		val = dst.Interface()
	} else {
		dst := reflect.New(typ).Elem()      //创建对象
		data := reflect.ValueOf(this.value) //源数据值
		dst.Set(data)                       //设置数据
		val = dst.Interface()
	}
	return &Var{value: val, kind: this.kind}
}

func (this *Var) Convert(target interface{}) error {
	kind := reflect.ValueOf(target).Kind()
	if kind != reflect.Ptr && kind != reflect.Map {
		return errors.New("Element.Bind must required a pointer or map argument")
	}

	switch this.kind {
	case reflect.String:
		return jsoniter.UnmarshalFromString(this.value.(string), target)
	default:
		bts, ok := this.value.([]byte)
		if ok {
			return jsoniter.Unmarshal(bts, target)
		}
		data, err := jsoniter.Marshal(this.value)
		if err != nil {
			return err
		}
		return jsoniter.Unmarshal(data, target)
	}
}

func (this Var) MarshalBinary() (data []byte, err error) {
	if this.isComparable {
		return this.Bytes(), nil
	} else {
		return jsoniter.Marshal(this.value)
	}
}

func (this *Var) UnmarshalBinary(data []byte) error {
	var value any
	if err := jsoniter.Unmarshal(data, &value); err != nil {
		return err
	}

	var elem = New(value)
	*this = *elem
	return nil
}

func (this Var) MarshalJSON() (data []byte, err error) {
	return this.MarshalBinary()
}

func (this *Var) UnmarshalJSON(data []byte) error {
	return this.UnmarshalBinary(data)
}

// Value return json value, implement driver.Valuer interface
func (v Var) Value() (driver.Value, error) {
	bytes, err := v.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}

// Scan value into Jsonb, implements sql.Scanner interface
func (this *Var) Scan(value interface{}) error {
	*this = *New(nil)
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

	return this.UnmarshalJSON(bytes)
}

// GormDataType gorm common data type
func (v Var) GormDataType() string {
	return v.kind.String()
}

// GormDBDataType gorm db data type
func (v Var) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		//return "JSON"
	case "mysql":
		if v.isComparable {
			switch v.kind {
			case reflect.String:
				return "varchar"
			case reflect.Int,
				reflect.Int8,
				reflect.Int16,
				reflect.Int32,
				reflect.Uint,
				reflect.Uint8,
				reflect.Uint16,
				reflect.Uint32:
				return "int"
			case reflect.Int64, reflect.Uint64:
				return "bigint"
			case reflect.Bool:
				return "tinyint"
			default:
				return ""
			}
		} else {
			return "JSON"
		}
	case "postgres":
		if v.isComparable {
			switch v.kind {
			case reflect.String:
				return "varchar"
			case reflect.Int,
				reflect.Int8,
				reflect.Int16,
				reflect.Int32,
				reflect.Uint,
				reflect.Uint8,
				reflect.Uint16,
				reflect.Uint32:
				return "integer"
			case reflect.Int64, reflect.Uint64:
				return "bigint"
			case reflect.Bool:
				return "bool"
			default:
				return ""
			}
		} else {
			return "JSONB"
		}
	case "sqlserver":
		//return "NVARCHAR(MAX)"
	}
	return ""
}

func (v Var) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if v.value == nil {
		return gorm.Expr("Null")
	}

	data, _ := v.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}

	return gorm.Expr("?", string(data))
}

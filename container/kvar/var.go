package kvar

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/golang-module/carbon/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/r3labs/diff/v3"
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
	isBytes      bool
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
		if _, ok := v.value.(time.Time); ok {
			v.isComparable = true
		}
		if _, ok := v.value.([]byte); ok {
			v.isBytes = true
		}
	}

	return v
}

func (this *Var) IsComparable() bool {
	return this.isComparable
}

func (this *Var) IsBytes() bool {
	return this.isBytes
}

func (this *Var) IsNil() bool {
	return this.value == nil
}

func (this *Var) IsNumeric() bool {
	return this.IsInt() || this.IsFloat() || this.IsUint()
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
	} else if this.IsBytes() || this.IsString() {
		return this.String() == ""
	} else if this.IsNumeric() {
		return this.Float64() == 0
	} else if this.IsSlice() {
		return len(this.Interfaces()) == 0
	} else if this.Bool() {
		return this.Bool() == false
	} else if this.IsMap() {
		return len(this.Map()) == 0
	} else if t, ok := this.value.(time.Time); ok {
		return t.IsZero()
	} else if this.IsStruct() || this.IsPointer() {
		return false
	} else {
		return true
	}
}

func (this *Var) String() string {
	if this.value == nil {
		return ""
	} else if this.isComparable {
		return cast.ToString(this.value)
	} else {
		str, _ := jsoniter.MarshalToString(this.value)
		return str
	}
}

func (this *Var) Strings() []string {
	vals := this.Interfaces()
	if vals == nil {
		return []string{}
	}

	strs := make([]string, 0)
	{
		for _, val := range vals {
			strs = append(strs, New(val).String())
		}
	}
	return strs

	//if this.IsSlice() {
	//	return cast.ToStringSlice(this.value)
	//} else if this.IsBytes() {
	//	strs := make([]string, 0)
	//	_ = jsoniter.Unmarshal(this.value.([]byte), &strs)
	//	return strs
	//} else if this.IsString() {
	//	strs := make([]string, 0)
	//	err := jsoniter.UnmarshalFromString(this.value.(string), &strs)
	//	if err == nil {
	//		return strs
	//	}
	//	return strings.Fields(this.value.(string))
	//}
	//return []string{}
}

func (this *Var) Int() int {
	if this.isComparable {
		return cast.ToInt(this.value)
	}
	return 0
}

func (this *Var) Ints() []int {
	vals := this.Interfaces()
	if vals == nil {
		return []int{}
	}

	ints := make([]int, 0)
	{
		for _, val := range vals {
			ints = append(ints, New(val).Int())
		}
	}
	return ints
}

func (this *Var) int8() int8 {
	if this.isComparable {
		return cast.ToInt8(this.value)
	}
	return 0
}

func (this *Var) Int8s() []int8 {
	vals := this.Interfaces()
	if vals == nil {
		return []int8{}
	}

	ints := make([]int8, 0)
	{
		for _, val := range vals {
			ints = append(ints, New(val).int8())
		}
	}
	return ints
}

func (this *Var) Int16() int16 {
	if this.isComparable {
		return cast.ToInt16(this.value)
	}
	return 0
}

func (this *Var) Int16s() []int16 {
	vals := this.Interfaces()
	if vals == nil {
		return []int16{}
	}

	ints := make([]int16, 0)
	{
		for _, val := range vals {
			ints = append(ints, New(val).Int16())
		}
	}
	return ints
}

func (this *Var) Int32() int32 {
	if this.isComparable {
		return cast.ToInt32(this.value)
	}
	return 0
}

func (this *Var) Int32s() []int32 {
	vals := this.Interfaces()
	if vals == nil {
		return []int32{}
	}

	ints := make([]int32, 0)
	{
		for _, val := range vals {
			ints = append(ints, New(val).Int32())
		}
	}
	return ints
}

func (this *Var) Int64() int64 {
	if this.isComparable {
		return cast.ToInt64(this.value)
	}
	return 0
}

func (this *Var) Int64s() []int64 {
	vals := this.Interfaces()
	if vals == nil {
		return []int64{}
	}

	ints := make([]int64, 0)
	{
		for _, val := range vals {
			ints = append(ints, New(val).Int64())
		}
	}
	return ints
}

func (this *Var) Float32() float32 {
	if this.isComparable {
		return cast.ToFloat32(this.value)
	}
	return 0
}

func (this *Var) Float32s() []float32 {
	vals := this.Interfaces()
	if vals == nil {
		return []float32{}
	}

	ints := make([]float32, 0)
	{
		for _, val := range vals {
			ints = append(ints, New(val).Float32())
		}
	}
	return ints
}

func (this *Var) Float64() float64 {
	return cast.ToFloat64(this.value)
}

func (this *Var) Float64s() []float64 {
	vals := this.Interfaces()
	if vals == nil {
		return []float64{}
	}

	ints := make([]float64, 0)
	{
		for _, val := range vals {
			ints = append(ints, New(val).Float64())
		}
	}
	return ints
}

func (this *Var) Uint() uint {
	return cast.ToUint(this.value)
}

func (this *Var) Uints() []uint {
	vals := this.Interfaces()
	if vals == nil {
		return []uint{}
	}

	ints := make([]uint, 0)
	{
		for _, val := range vals {
			ints = append(ints, New(val).Uint())
		}
	}
	return ints
}

func (this *Var) Uint8() uint8 {
	return cast.ToUint8(this.value)
}

func (this *Var) Uint8s() []uint8 {
	vals := this.Interfaces()
	if vals == nil {
		return []uint8{}
	}

	ints := make([]uint8, 0)
	{
		for _, val := range vals {
			ints = append(ints, New(val).Uint8())
		}
	}
	return ints
}

func (this *Var) Uint16() uint16 {
	return cast.ToUint16(this.value)
}

func (this *Var) Uint16s() []uint16 {
	vals := this.Interfaces()
	if vals == nil {
		return []uint16{}
	}

	ints := make([]uint16, 0)
	{
		for _, val := range vals {
			ints = append(ints, New(val).Uint16())
		}
	}
	return ints
}

func (this *Var) Uint32() uint32 {
	return cast.ToUint32(this.value)
}

func (this *Var) Uint32s() []uint32 {
	vals := this.Interfaces()
	if vals == nil {
		return []uint32{}
	}

	ints := make([]uint32, 0)
	{
		for _, val := range vals {
			ints = append(ints, New(val).Uint32())
		}
	}
	return ints
}

func (this *Var) Uint64() uint64 {
	return cast.ToUint64(this.value)
}

func (this *Var) Uint64s() []uint64 {
	vals := this.Interfaces()
	if vals == nil {
		return []uint64{}
	}

	ints := make([]uint64, 0)
	{
		for _, val := range vals {
			ints = append(ints, New(val).Uint64())
		}
	}
	return ints
}

func (this *Var) Bool() bool {
	return cast.ToBool(this.String())
}

func (this *Var) Bools() []bool {
	vals := this.Interfaces()
	if vals == nil {
		return []bool{}
	}

	fmt.Println("vals", vals)

	ints := make([]bool, 0)
	{
		for _, val := range vals {
			ints = append(ints, New(val).Bool())
		}
	}
	return ints
}

func (this *Var) Interface() interface{} {
	return this.value
}

func (this *Var) Interfaces() []interface{} {
	if this.value == nil {
		return []interface{}{}
	}

	if this.IsBytes() {
		strs := make([]interface{}, 0)
		_ = jsoniter.Unmarshal(this.value.([]byte), &strs)
		return strs
	} else if this.IsString() {
		vals := make([]interface{}, 0)
		err := jsoniter.UnmarshalFromString(this.value.(string), &vals)
		if err == nil {
			return vals
		}

		strs := strings.Fields(this.value.(string))
		if strs != nil && len(strs) > 0 {
			for _, str := range strs {
				vals = append(vals, str)
			}
		}
		return vals
	} else {
		arr := make([]interface{}, 0)
		btes, err := jsoniter.Marshal(this.value)
		if err != nil {
			return arr
		}

		_ = jsoniter.Unmarshal(btes, &arr)
		return arr
	}
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

	if t, ok := this.value.(time.Time); ok {
		return t
	} else if t, err := parseTimestamp(); err == nil {
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
	if this.value == nil {
		return []byte{}
	} else if this.IsBytes() {
		return this.value.([]byte)
	} else if this.isComparable {
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
	df, _ := diff.Diff(this.value, elem.value)
	return len(df) == 0
}

func (this *Var) Map() map[string]interface{} {
	if this.value == nil {
		return map[string]interface{}{}
	} else if this.IsMap() {
		return cast.ToStringMap(this.value)
	} else if this.IsBytes() {
		mp := make(map[string]interface{})
		_ = jsoniter.Unmarshal(this.value.([]byte), &mp)
		return mp
	} else if this.IsStruct() || this.IsPointer() {
		bts, err := jsoniter.Marshal(this.Interface())
		if err != nil {
			return map[string]interface{}{}
		}

		mp := make(map[string]interface{})
		_ = jsoniter.Unmarshal(bts, &mp)
		return mp
	} else if this.IsString() {
		mp := make(map[string]interface{})
		_ = jsoniter.UnmarshalFromString(this.value.(string), &mp)
		return mp
	} else {
		return cast.ToStringMap(this.value)
	}
}

func (this *Var) MapStr() map[string]string {
	mp := this.Map()
	if mp == nil {
		return map[string]string{}
	}

	result := make(map[string]string)
	for k, v := range mp {
		mp[k] = New(v).String()
	}
	return result
}

func (this *Var) MapInt() map[string]int {
	mp := this.Map()
	if mp == nil {
		return map[string]int{}
	}

	result := make(map[string]int)
	for k, v := range mp {
		mp[k] = New(v).Int()
	}
	return result
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
	mp := this.MapInt()
	if mp == nil {
		return map[string]int64{}
	}
	newMp := make(map[string]int64)
	{
		for s, i := range mp {
			newMp[s] = cast.ToInt64(i)
		}
	}
	return newMp
}

func (this *Var) MapFloat32() map[string]float32 {
	mp := this.MapInt64()
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
	mp := this.MapInt64()
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
	mp := this.Map()
	if mp == nil {
		return map[string]bool{}
	}
	newMp := make(map[string]bool)
	{
		for s, i := range mp {
			newMp[s] = New(i).Bool()
		}
	}
	return cast.ToStringMapBool(this.value)
}

func (this *Var) MapVar() map[string]*Var {
	mp := this.Map()
	if mp == nil || len(mp) == 0 {
		return map[string]*Var{}
	}
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
	if this.value == nil {
		return nil
	}
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
		return "JSON"
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
		return "NVARCHAR(MAX)"
	}
	return ""
}

func (v Var) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if v.value == nil {
		return gorm.Expr("Null")
	}

	if v.isComparable {
		return gorm.Expr("?", v.value)
	}

	data, _ := v.MarshalJSON()
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

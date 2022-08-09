package kvar

import (
	"context"
	"database/sql/driver"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	type args struct {
		val any
	}
	tests := []struct {
		name string
		args args
		want *Var
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.val); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Bool(t *testing.T) {
	type fields struct {
		val *Var
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "string - true",
			fields: fields{New("true")},
			want:   true,
		},
		{
			name:   "string - false",
			fields: fields{New("false")},
			want:   false,
		},
		{
			name:   "string - TRUE",
			fields: fields{New("TRUE")},
			want:   true,
		},
		{
			name:   "string - FALSE",
			fields: fields{New("FALSE")},
			want:   false,
		},
		{
			name:   "string - True",
			fields: fields{New("True")},
			want:   true,
		},
		{
			name:   "string - False",
			fields: fields{New("False")},
			want:   false,
		},
		{
			name:   "int - 1",
			fields: fields{New(1)},
			want:   true,
		},
		{
			name:   "int - 0",
			fields: fields{New(0)},
			want:   false,
		},
		{
			name:   "int - -1",
			fields: fields{New(-1)},
			want:   true,
		},
		{
			name:   "int - 2",
			fields: fields{New(2)},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := tt.fields.val
			if got := this.Bool(); got != tt.want {
				t.Errorf("Bool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Bools(t *testing.T) {
	type fields struct {
		val *Var
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "strings",
			fields: fields{New([]string{"true", "false", "True", "False", "TRUE", "FALSE"})},
		},
		{
			name:   "int",
			fields: fields{New([]int{1, 0, 2, -1})},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := tt.fields.val
			got := this.Bools()
			t.Log(tt.name, got)
		})
	}
}

func TestVar_Bytes(t *testing.T) {
	type person struct {
		Name string `json:"name"`
	}
	type fields struct {
		value *Var
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name:   "nil",
			fields: fields{New(nil)},
			want:   []byte{},
		},
		{
			name:   "bytes",
			fields: fields{New([]byte("ssss"))},
			want:   []byte("ssss"),
		},
		{
			name:   "string",
			fields: fields{New("ssss")},
			want:   []byte("ssss"),
		},
		{
			name:   "int",
			fields: fields{New(11)},
			want:   []byte(fmt.Sprintf("%d", 11)),
		},
		{
			name:   "array",
			fields: fields{New([]string{"a", "b"})},
			want:   []byte(`["a","b"]`),
		},
		{
			name:   "map",
			fields: fields{New(map[string]any{"a": "aa", "b": 1})},
			want:   []byte(`{"a":"aa","b":1}`),
		},
		{
			name:   "arrayMap",
			fields: fields{New([]map[string]any{{"a": "aa"}})},
			want:   []byte(`[{"a":"aa"}]`),
		},
		{
			name:   "struct",
			fields: fields{New(person{Name: "ss"})},
			want:   []byte(`{"name":"ss"}`),
		},
		{
			name:   "pointer",
			fields: fields{New(&person{Name: "ss"})},
			want:   []byte(`{"name":"ss"}`),
		},
		{
			name:   "boolean",
			fields: fields{New(false)},
			want:   []byte("false"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := tt.fields.value
			if got := this.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bytes() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestVar_Clone(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   *Var
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Convert(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	type args struct {
		target interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if err := this.Convert(tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVar_Duration(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Duration(); got != tt.want {
				t.Errorf("Duration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Equal(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	type args struct {
		elem *Var
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Equal(tt.args.elem); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Float32(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   float32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Float32(); got != tt.want {
				t.Errorf("Float32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Float32s(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []float32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Float32s(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Float32s() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Float64(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Float64(); got != tt.want {
				t.Errorf("Float64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Float64s(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Float64s(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Float64s() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_GormDBDataType(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	type args struct {
		db    *gorm.DB
		field *schema.Field
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := v.GormDBDataType(tt.args.db, tt.args.field); got != tt.want {
				t.Errorf("GormDBDataType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_GormDataType(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := v.GormDataType(); got != tt.want {
				t.Errorf("GormDataType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_GormValue(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	type args struct {
		ctx context.Context
		db  *gorm.DB
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   clause.Expr
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := v.GormValue(tt.args.ctx, tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GormValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Int(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Int(); got != tt.want {
				t.Errorf("Int() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Int16(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   int16
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Int16(); got != tt.want {
				t.Errorf("Int16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Int16s(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []int16
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Int16s(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int16s() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Int32(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   int32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Int32(); got != tt.want {
				t.Errorf("Int32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Int32s(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []int32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Int32s(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int32s() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Int64(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Int64(); got != tt.want {
				t.Errorf("Int64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Int64s(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Int64s(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int64s() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Int8s(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []int8
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Int8s(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int8s() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Interface(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Interface(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Interface() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Interfaces(t *testing.T) {
	jsonStr := `["a","b",1,true]`
	bts := []byte(jsonStr)
	type fields struct {
		v *Var
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "anys",
			fields: fields{New([]any{1, "a", 1.2, true, map[string]any{"aa": 1, "bb": "cc"}})},
		},
		{
			name:   "strings",
			fields: fields{New([]string{"a", "b"})},
		},
		{
			name:   "string",
			fields: fields{New("abcd")},
		},
		{
			name:   "json string",
			fields: fields{New(jsonStr)},
		},
		{
			name:   "json bytes",
			fields: fields{New(bts)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := tt.fields.v
			got := this.Interfaces()
			if got == nil || len(got) == 0 {
				t.Errorf("Interfaces() = %v", got)
				return
			}

			t.Log(tt.name, got)
		})
	}
}

func TestVar_Ints(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Ints(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_IsBytes(t *testing.T) {
	type fields struct {
		value *Var
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "bytes",
			fields: fields{New([]byte("ssss"))},
			want:   true,
		},
		{
			name:   "string",
			fields: fields{New("sss")},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := tt.fields.value
			if got := this.IsBytes(); got != tt.want {
				t.Errorf("IsBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_IsComparable(t *testing.T) {
	type person struct {
		Name string `json:"name"`
	}
	type fields struct {
		value *Var
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "string",
			fields: fields{New("aaa")},
			want:   true,
		},
		{
			name:   "numbric",
			fields: fields{New(111)},
			want:   true,
		},
		{
			name:   "json string",
			fields: fields{New(`{}`)},
			want:   true,
		},
		{
			name:   "json bytes",
			fields: fields{New([]byte(`{}`))},
			want:   false,
		},
		{
			name:   "array",
			fields: fields{New([]string{"a", "b"})},
			want:   false,
		},
		{
			name:   "map",
			fields: fields{New(map[string]any{"aa": "bb"})},
			want:   false,
		},
		{
			name:   "struct",
			fields: fields{New(person{Name: "kit"})},
			want:   false,
		},
		{
			name:   "pointer",
			fields: fields{New(new(person))},
			want:   false,
		},
		{
			name:   "boolean",
			fields: fields{New(false)},
			want:   true,
		},
		{
			name:   "time",
			fields: fields{New(time.Now())},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := tt.fields.value
			if got := this.IsComparable(); got != tt.want {
				t.Errorf("IsComparable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_IsEmpty(t *testing.T) {
	type person struct {
		Name string `json:"name"`
	}
	type fields struct {
		value *Var
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "any - nil",
			fields: fields{New(nil)},
			want:   true,
		},
		{
			name:   "any - not nil",
			fields: fields{New(any(333))},
			want:   false,
		},
		{
			name:   "string - empty",
			fields: fields{New("")},
			want:   true,
		},
		{
			name:   "string - not empty",
			fields: fields{New("aa")},
			want:   false,
		},
		{
			name:   "number - empty",
			fields: fields{New(0)},
			want:   true,
		},
		{
			name:   "number - not empty",
			fields: fields{New(-1)},
			want:   false,
		},
		{
			name:   "boolean - true",
			fields: fields{New(true)},
			want:   false,
		},
		{
			name:   "boolean - false",
			fields: fields{New(false)},
			want:   true,
		},
		{
			name:   "time - not empty",
			fields: fields{New(time.Now())},
			want:   false,
		},
		{
			name:   "time - empty",
			fields: fields{New(time.Time{})},
			want:   true,
		},
		{
			name:   "array - empty",
			fields: fields{New([]string{})},
			want:   true,
		},
		{
			name:   "array - not empty",
			fields: fields{New([]string{"a", "b"})},
			want:   false,
		},
		{
			name:   "array - all nil",
			fields: fields{New([]any{nil, nil})},
			want:   false,
		},
		{
			name:   "map - empty",
			fields: fields{New(map[string]any{})},
			want:   true,
		},
		{
			name:   "map - not empty",
			fields: fields{New(map[string]any{"a": 1})},
			want:   false,
		},
		{
			name:   "map - has nil",
			fields: fields{New(map[string]any{"a": nil})},
			want:   false,
		},
		{
			name:   "struct - empty",
			fields: fields{New(person{})},
			want:   false,
		},
		{
			name:   "struct - not empty",
			fields: fields{New(person{Name: "sss"})},
			want:   false,
		},
		{
			name:   "pointer - empty",
			fields: fields{New(new(person))},
			want:   false,
		},
		{
			name:   "pointer - not empty",
			fields: fields{New(&person{Name: "sss"})},
			want:   false,
		},
		{
			name:   "pointer - nil",
			fields: fields{New((*person)(nil))},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := tt.fields.value
			if got := this.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_IsFloat(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.IsFloat(); got != tt.want {
				t.Errorf("IsFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_IsInt(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.IsInt(); got != tt.want {
				t.Errorf("IsInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_IsMap(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.IsMap(); got != tt.want {
				t.Errorf("IsMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_IsNil(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.IsNil(); got != tt.want {
				t.Errorf("IsNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_IsNumeric(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.IsNumeric(); got != tt.want {
				t.Errorf("IsNumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_IsPointer(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.IsPointer(); got != tt.want {
				t.Errorf("IsPointer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_IsSlice(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.IsSlice(); got != tt.want {
				t.Errorf("IsSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_IsString(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.IsString(); got != tt.want {
				t.Errorf("IsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_IsStruct(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.IsStruct(); got != tt.want {
				t.Errorf("IsStruct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_IsUint(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.IsUint(); got != tt.want {
				t.Errorf("IsUint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Kind(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   reflect.Kind
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Kind(); got != tt.want {
				t.Errorf("Kind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Map(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Map(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_MapBool(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.MapBool(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_MapFloat32(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]float32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.MapFloat32(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapFloat32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_MapFloat64(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.MapFloat64(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_MapInt(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.MapInt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_MapInt16(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]int16
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.MapInt16(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapInt16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_MapInt32(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]int32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.MapInt32(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_MapInt64(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.MapInt64(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_MapInt8(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]int8
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.MapInt8(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapInt8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_MapSlice(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string][]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.MapSlice(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_MapSliceStrings(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string][]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.MapSliceStrings(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapSliceStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_MapSliceVars(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string][]*Var
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.MapSliceVars(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapSliceVars() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_MapStr(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.MapStr(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_MapVar(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]*Var
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.MapVar(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapVar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_MarshalBinary(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name     string
		fields   fields
		wantData []byte
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			gotData, err := this.MarshalBinary()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("MarshalBinary() gotData = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

func TestVar_MarshalJSON(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name     string
		fields   fields
		wantData []byte
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			gotData, err := this.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("MarshalJSON() gotData = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

func TestVar_Scan(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if err := this.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVar_String(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Strings(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Strings(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Strings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Time(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Time(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Times(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Times(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Times() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Uint(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   uint
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Uint(); got != tt.want {
				t.Errorf("Uint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Uint16(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   uint16
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Uint16(); got != tt.want {
				t.Errorf("Uint16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Uint16s(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []uint16
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Uint16s(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uint16s() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Uint32(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Uint32(); got != tt.want {
				t.Errorf("Uint32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Uint32s(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Uint32s(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uint32s() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Uint64(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Uint64(); got != tt.want {
				t.Errorf("Uint64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Uint64s(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Uint64s(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uint64s() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Uint8(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   uint8
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Uint8(); got != tt.want {
				t.Errorf("Uint8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Uint8s(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []uint8
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Uint8s(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uint8s() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Uints(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []uint
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Uints(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_UnmarshalBinary(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if err := this.UnmarshalBinary(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVar_UnmarshalJSON(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if err := this.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVar_Value(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    driver.Value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			got, err := v.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Value() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_Vars(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []*Var
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.Vars(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Vars() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_int8(t *testing.T) {
	type fields struct {
		value        any
		kind         reflect.Kind
		isComparable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   int8
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Var{
				value:        tt.fields.value,
				kind:         tt.fields.kind,
				isComparable: tt.fields.isComparable,
			}
			if got := this.int8(); got != tt.want {
				t.Errorf("int8() = %v, want %v", got, tt.want)
			}
		})
	}
}

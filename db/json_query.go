package db

//
//import (
//	"fmt"
//	"github.com/gogf/gf/os/glog"
//	jsoniter "github.com/json-iterator/go"
//	"github.com/spf13/cast"
//	"gorm.io/gorm"
//	"gorm.io/gorm/clause"
//	"reflect"
//)
//
//type JsonQuery struct {
//	column   string
//	multi    bool
//	currFunc string
//	keys     []string
//	values   []interface{}
//}
//
//func NewJsonQuery(column string) *JsonQuery {
//	return &JsonQuery{column: column}
//}
//
//func (this *JsonQuery) Build(builder clause.Builder) {
//	glog.Info("key", this.keys)
//	stmt, ok := builder.(*gorm.Statement)
//	if !ok {
//		return
//	}
//
//	keys := "$"
//	if this.keys != nil && len(this.keys) > 0 {
//		for _, v := range this.keys {
//			num, err := cast.ToInt64E(v)
//			if err == nil {
//				keys += fmt.Sprintf("[%d]", num)
//			} else {
//				keys += fmt.Sprintf(".%s", v)
//			}
//		}
//	}
//
//	switch this.currFunc {
//	case "HasKey":
//		if len(this.keys) > 0 {
//			_, _ = builder.WriteString(fmt.Sprintf("JSON_EXTRACT(%s, '%s') IS NOT NULL", stmt.Quote(this.column), keys))
//		}
//	case "Eq", "In":
//		if this.values == nil || len(this.values) == 0 {
//			break
//		}
//
//		if this.multi {
//			vals := make([]interface{}, 0)
//			for _, v := range this.values {
//				switch reflect.ValueOf(v).Kind() {
//				case reflect.String:
//					vals = append(vals, fmt.Sprintf("\"%s\"", v.(string)))
//					//vals = append(vals, v.(string))
//				default:
//					vals = append(vals, fmt.Sprintf("%v", v))
//				}
//			}
//
//			relation := " AND "
//			if this.currFunc == "In" {
//				relation = " OR "
//			}
//			_, _ = builder.WriteString("(")
//			for idx, v := range vals {
//				//_, _ = builder.WriteString(fmt.Sprintf("JSON_CONTAINS(%s, '", stmt.Quote(this.column)))
//				_, _ = builder.WriteString(fmt.Sprintf("JSON_CONTAINS(%s, ", stmt.Quote(this.column)))
//				stmt.AddVar(builder, v)
//				//_, _ = builder.WriteString(fmt.Sprintf("', '%s')", keys))
//				_, _ = builder.WriteString(fmt.Sprintf(", '%s')", keys))
//				if idx < len(vals)-1 {
//					_, _ = builder.WriteString(relation)
//				}
//			}
//			_, _ = builder.WriteString(")")
//		} else {
//			if this.currFunc == "Eq" {
//				_, _ = builder.WriteString(fmt.Sprintf("JSON_EXTRACT(%s, '%s') = ", stmt.Quote(this.column), keys))
//				stmt.AddVar(builder, this.values[0])
//			} else {
//				_, _ = builder.WriteString(fmt.Sprintf("JSON_EXTRACT(%s, '%s') IN ", stmt.Quote(this.column), keys))
//				stmt.AddVar(builder, this.values)
//			}
//		}
//	case "Neq", "NotIn":
//		if this.values == nil || len(this.values) == 0 {
//			break
//		}
//
//		if this.multi {
//			vals := make([]interface{}, 0)
//			for _, v := range this.values {
//				switch reflect.ValueOf(v).Kind() {
//				case reflect.String:
//					vals = append(vals, fmt.Sprintf("\"%s\"", v.(string)))
//				default:
//					vals = append(vals, v)
//				}
//			}
//
//			relation := " AND "
//			//if this.currFunc == "In" {
//			//    relation = " OR "
//			//}
//			_, _ = builder.WriteString("(")
//			for idx, v := range vals {
//				_, _ = builder.WriteString(fmt.Sprintf("JSON_CONTAINS(%s, '", stmt.Quote(this.column)))
//				stmt.AddVar(builder, v)
//				_, _ = builder.WriteString(fmt.Sprintf("', '%s')", keys))
//				if idx < len(vals)-1 {
//					_, _ = builder.WriteString(relation)
//				}
//			}
//			_, _ = builder.WriteString(")")
//		} else {
//			if this.currFunc == "Neq" {
//				_, _ = builder.WriteString(fmt.Sprintf("JSON_EXTRACT(%s, '%s') != ", stmt.Quote(this.column), keys))
//				stmt.AddVar(builder, this.values[0])
//			} else {
//				_, _ = builder.WriteString(fmt.Sprintf("JSON_EXTRACT(%s, '%s') NOT IN ", stmt.Quote(this.column), keys))
//				stmt.AddVar(builder, this.values)
//			}
//		}
//	case "Gt":
//		if this.values == nil || len(this.values) == 0 {
//			break
//		}
//		_, _ = builder.WriteString(fmt.Sprintf("%s->>'%s' > ", this.column, keys))
//		stmt.AddVar(builder, this.values[0])
//	case "Gte":
//		if this.values == nil || len(this.values) == 0 {
//			break
//		}
//		_, _ = builder.WriteString(fmt.Sprintf("%s->>'%s' >= ", this.column, keys))
//		stmt.AddVar(builder, this.values[0])
//	case "Lte":
//		if this.values == nil || len(this.values) == 0 {
//			break
//		}
//		_, _ = builder.WriteString(fmt.Sprintf("%s->>'%s' <= ", this.column, keys))
//		stmt.AddVar(builder, this.values[0])
//	case "Lt":
//		if this.values == nil || len(this.values) == 0 {
//			break
//		}
//		_, _ = builder.WriteString(fmt.Sprintf("%s->>'%s' < ", this.column, keys))
//		stmt.AddVar(builder, this.values[0])
//	case "Between":
//		if this.values == nil || len(this.values) != 2 {
//			break
//		}
//		if this.multi {
//			_, _ = builder.WriteString(fmt.Sprintf("(%s->>'%s[0]' >= ", this.column, keys))
//			stmt.AddVar(builder, this.values[0])
//			_, _ = builder.WriteString(fmt.Sprintf(" AND %s->>'%s[1]' <= ", this.column, keys))
//			stmt.AddVar(builder, this.values[1])
//			_, _ = builder.WriteString(")")
//		} else {
//			_, _ = builder.WriteString(fmt.Sprintf("%s->>'%s' BETWEEN ", this.column, keys))
//			stmt.AddVar(builder, this.values[0])
//			_, _ = builder.WriteString(" AND ")
//			stmt.AddVar(builder, this.values[1])
//			_, _ = builder.WriteString(")")
//		}
//		//return fmt.Sprintf("(%s->>'%s' BETWEEN ? AND ?)", this.column, keys), this.values, nil
//	case "Like":
//		_, _ = builder.WriteString(fmt.Sprintf("JSON_SEARCH(%s, 'one', ", this.column))
//		stmt.AddVar(builder, this.values[0])
//		_, _ = builder.WriteString(fmt.Sprintf(", NULL, '%s') IS NOT NULL", keys))
//	case "IsNull":
//		_, _ = builder.WriteString(fmt.Sprintf("%s->>'%s' IS NULL", this.column, keys))
//	case "NotNull":
//		_, _ = builder.WriteString(fmt.Sprintf("%s->>'%s' IS NOT NULL", this.column, keys))
//	default:
//		return
//	}
//}
//
////func (this *JsonQuery) ToSql() (string, []interface{}, error) {
////    keys := "$"
////    if this.keys != nil && len(this.keys) > 0 {
////        for _, v := range this.keys {
////            num, err := cast.ToInt64E(v)
////            if err == nil {
////                keys += fmt.Sprintf("[%d]", num)
////            } else {
////                keys += fmt.Sprintf(".%s", v)
////            }
////        }
////    }
////
////    switch this.currFunc {
////    case "HasKey":
////        if len(this.keys) > 0 {
////            return fmt.Sprintf("JSON_EXTRACT(%s, '%s') IS NOT NULL", this.column, keys), nil, nil
////        }
////    case "Eq", "In":
////        if this.values == nil || len(this.values) == 0 {
////            return "", nil, nil
////        }
////
////        if this.multi {
////
////        } else {
////            if this.currFunc == "Eq" {
////                return fmt.Sprintf("JSON_EXTRACT(%s, '%s') = ?", this.column, keys), this.values, nil
////            } else {
////                return fmt.Sprintf("JSON_EXTRACT(%s, '%s') IN ?", this.column, keys), this.values, nil
////            }
////        }
////    }
////    return "", nil, nil
////}
//
////func (this *JsonQuery) ParseDataframe(df dataframe.DataFrame) dataframe.DataFrame {
////    return dataframe.DataFrame{}
////}
//
//func (this *JsonQuery) HasKey(keys []string) *JsonQuery {
//	this.currFunc = "HasKey"
//	this.keys = keys
//	return this
//}
//
//func (this *JsonQuery) IsNull(keys []string) *JsonQuery {
//	this.keys = keys
//	this.currFunc = "IsNull"
//	return this
//}
//
//func (this *JsonQuery) NotNull(keys []string) *JsonQuery {
//	this.keys = keys
//	this.currFunc = "NotNull"
//	return this
//}
//
//func (this *JsonQuery) Eq(value interface{}, keys []string, keysIsArray bool) *JsonQuery {
//	this.multi = keysIsArray
//	this.values = make([]interface{}, 0)
//	this.currFunc = "Eq"
//	this.parseValue(value)
//	this.keys = keys
//	return this
//}
//
//func (this *JsonQuery) Neq(value interface{}, keys []string, keysIsArray bool) *JsonQuery {
//	this.multi = keysIsArray
//	this.values = make([]interface{}, 0)
//	this.currFunc = "Neq"
//	this.parseValue(value)
//	this.keys = keys
//	return this
//}
//
//func (this *JsonQuery) Gt(value interface{}, keys []string) *JsonQuery {
//	this.currFunc = "Gt"
//	this.keys = keys
//	this.parseValue(value)
//	return this
//}
//
//func (this *JsonQuery) Gte(value interface{}, keys []string) *JsonQuery {
//	this.currFunc = "Gte"
//	this.keys = keys
//	this.parseValue(value)
//	return this
//}
//
//func (this *JsonQuery) Lt(value interface{}, keys []string) *JsonQuery {
//	this.currFunc = "Lt"
//	this.keys = keys
//	this.parseValue(value)
//	return this
//}
//
//func (this *JsonQuery) Lte(value interface{}, keys []string) *JsonQuery {
//	this.currFunc = "Lte"
//	this.keys = keys
//	this.parseValue(value)
//	return this
//}
//
//func (this *JsonQuery) Between(value interface{}, keys []string, keysIsArray bool) *JsonQuery {
//	this.parseValue(value)
//	if len(this.values) != 2 {
//		this.currFunc = ""
//		return this
//	}
//
//	this.multi = keysIsArray
//	this.currFunc = "Between"
//	this.keys = keys
//	return this
//}
//
//func (this *JsonQuery) In(value interface{}, keys []string, keysIsArray bool) *JsonQuery {
//	this.multi = keysIsArray
//	this.currFunc = "In"
//	this.parseValue(value)
//	this.keys = keys
//	return this
//}
//
//func (this *JsonQuery) NotIn(value interface{}, keys []string, keysIsArray bool) *JsonQuery {
//	this.multi = keysIsArray
//	this.currFunc = "NotIn"
//	this.parseValue(value)
//	this.keys = keys
//	return this
//}
//
//func (this *JsonQuery) Like(value interface{}, keys []string) *JsonQuery {
//	this.currFunc = "Like"
//	this.parseValue(value)
//	this.keys = keys
//	return this
//}
//
//func (this *JsonQuery) parseValue(value interface{}) {
//	this.values = make([]interface{}, 0)
//	if value != nil {
//		reflectValue := reflect.Indirect(reflect.ValueOf(value))
//		switch reflectValue.Kind() {
//		case reflect.Slice, reflect.Array:
//			if data, err := jsoniter.Marshal(value); err == nil {
//				_ = jsoniter.Unmarshal(data, &this.values)
//			}
//		default:
//			this.values = append(this.values, value)
//		}
//	}
//}

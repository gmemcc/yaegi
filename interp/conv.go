package interp

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/spf13/cast"
	"go/constant"
	"reflect"
	"strconv"
	"strings"
)

func canIconv(typ *itype, expected *itype) bool {
	if typ.assignableTo(expected) {
		return true
	}
	if typ.rtype.Kind() == reflect.String {
		return true
	}
	ertype := expected.rtype
	_, err := rconv(reflect.New(typ.rtype).Elem(), ertype)
	return err == nil
}

func canIconvBool(t *itype) bool {
	typ := t.TypeOf()
	return typ.Kind() == reflect.Bool || isNumber(typ) || isString(typ) || isInterface(t)
}

func canRconvBool(t reflect.Type) bool {
	return t.Kind() == reflect.Bool || t.Kind() == reflect.Interface || isNumber(t) || isString(t)
}

func rconv(src reflect.Value, expectedType reflect.Type) (reflect.Value, error) {
	if !src.IsValid() {
		return src, nil
	}
	srcType := src.Type()
	if srcType == expectedType {
		return src, nil
	}
	if srcType.Kind() == expectedType.Kind() && srcType.Kind() != reflect.Struct &&
		(srcType.PkgPath() != expectedType.PkgPath() || srcType.Name() != expectedType.Name()) {
		// type def from existing type
		return src.Convert(expectedType), nil
	}
	value := src.Interface()
	switch expectedType.Kind() {
	case reflect.Bool:
		casted, err := cast.ToBoolE(value)
		if err == nil {
			return reflect.ValueOf(casted), nil
		} else {
			return src, err
		}
	case reflect.Int:
		casted, err := cast.ToIntE(value)
		if err == nil {
			return reflect.ValueOf(casted), nil
		} else {
			return src, err
		}
	case reflect.Int8:
		casted, err := cast.ToInt8E(value)
		if err == nil {
			return reflect.ValueOf(casted), nil
		} else {
			return src, err
		}
	case reflect.Int16:
		casted, err := cast.ToInt16E(value)
		if err == nil {
			return reflect.ValueOf(casted), nil
		} else {
			return src, err
		}
	case reflect.Int32:
		casted, err := cast.ToInt32E(value)
		if err == nil {
			return reflect.ValueOf(casted), nil
		} else {
			return src, err
		}
	case reflect.Int64:
		casted, err := cast.ToInt64E(value)
		if err == nil {
			return reflect.ValueOf(casted), nil
		} else {
			return src, err
		}
	case reflect.Uint:
		casted, err := cast.ToUintE(value)
		if err == nil {
			return reflect.ValueOf(casted), nil
		} else {
			return src, err
		}
	case reflect.Uint8:
		casted, err := cast.ToUint8E(value)
		if err == nil {
			return reflect.ValueOf(casted), nil
		} else {
			return src, err
		}
	case reflect.Uint16:
		casted, err := cast.ToUint16E(value)
		if err == nil {
			return reflect.ValueOf(casted), nil
		} else {
			return src, err
		}
	case reflect.Uint32:
		casted, err := cast.ToUint32E(value)
		if err == nil {
			return reflect.ValueOf(casted), nil
		} else {
			return src, err
		}
	case reflect.Uint64:
		casted, err := cast.ToUint64E(value)
		if err == nil {
			return reflect.ValueOf(casted), nil
		} else {
			return src, err
		}
	case reflect.Float32:
		casted, err := cast.ToFloat32E(value)
		if err == nil {
			return reflect.ValueOf(casted), nil
		} else {
			return src, err
		}
	case reflect.Float64:
		casted, err := cast.ToFloat64E(value)
		if err == nil {
			return reflect.ValueOf(casted), nil
		} else {
			return src, err
		}
	case reflect.String:
		switch reflect.ValueOf(value).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			casted := fmt.Sprintf("%d", value)
			return reflect.ValueOf(casted), nil
		case reflect.Map:
			bytes, err := json.Marshal(value)
			return reflect.ValueOf(string(bytes)), err
		default:
			casted, err := cast.ToStringE(value)
			if err == nil {
				return reflect.ValueOf(casted), nil
			} else {
				return src, err
			}
		}
	case reflect.Struct:
		castedPtrValue := reflect.New(expectedType)
		indirect := reflect.Indirect(src)
		kind := indirect.Kind()
		switch kind {
		case reflect.String:
			// assume value is in json format
			err := json.Unmarshal([]byte(indirect.String()), castedPtrValue.Interface())
			if err == nil {
				return castedPtrValue.Elem(), nil
			} else {
				return src, err
			}
		case reflect.Struct:
			err := copier.Copy(castedPtrValue.Interface(), value)
			if err == nil {
				return castedPtrValue.Elem(), nil
			} else {
				return src, err
			}
		case reflect.Map:
			var bytes []byte
			var err error
			bytes, err = json.Marshal(indirect.Interface())
			if err == nil {
				err = json.Unmarshal(bytes, castedPtrValue.Interface())
				if err == nil {
					return castedPtrValue.Elem(), nil
				}
			}
			return src, err
		case reflect.Interface:
			return indirect.Elem().Convert(expectedType), nil
		default:
			return src, errors.New(fmt.Sprintf(""))
		}
	case reflect.Map:
		indirect := rconvToConcrete(reflect.Indirect(src))
		kind := indirect.Kind()
		switch kind {
		case reflect.String:
			castedPtrValue := reflect.New(expectedType)
			str := indirect.String()
			if str == "" {
				return castedPtrValue.Elem(), nil
			}
			err := json.Unmarshal([]byte(str), castedPtrValue.Interface())
			if err == nil {
				return castedPtrValue.Elem(), nil
			} else {
				return src, err
			}
		case reflect.Map:
			break
		case reflect.Invalid:
			return reflect.Zero(expectedType), nil
		case reflect.Struct:
			srcBytes, err := json.Marshal(src.Interface())
			if err != nil {
				return src, err
			} else {
				castedPtrValue := reflect.New(expectedType)
				err := json.Unmarshal(srcBytes, castedPtrValue.Interface())
				if err == nil {
					return castedPtrValue.Elem(), nil
				} else {
					return src, err
				}
			}
		default:
			return src, nil
		}
		ktype := expectedType.Key()
		vtype := expectedType.Elem()
		if ktype == indirect.Type().Key() && vtype == indirect.Type().Elem() {
			return indirect, nil
		}
		castedValue := reflect.MakeMapWithSize(expectedType, 0)
		keys := indirect.MapKeys()
		for i := 0; i < len(keys); i++ {
			k := keys[i]
			v := indirect.MapIndex(k)
			var kcasted, vcasted reflect.Value
			var err error
			kcasted, err = rconv(k, ktype)
			if err != nil {
				return src, err
			}
			vcasted, err = rconv(v, vtype)
			if err != nil {
				return src, err
			}
			castedValue.SetMapIndex(kcasted, vcasted)
		}
		return castedValue, nil
	case reflect.Slice:
		src = rconvToConcrete(src)
		srcType = src.Type()
		if srcType.Kind() == reflect.String {
			return rconv(reflect.ValueOf([]uint8(src.String())), expectedType)
		} else if srcType.Kind() != reflect.Slice {
			return src, nil
		}
		vtype := expectedType.Elem()
		castedValue := reflect.MakeSlice(expectedType, src.Len(), src.Cap())
		for i := 0; i < src.Len(); i++ {
			vcasted, err := rconv(src.Index(i), vtype)
			if err == nil {
				castedValue.Index(i).Set(vcasted)
			} else {
				return src, err
			}
		}
		return castedValue, nil
	case reflect.Ptr:
		castedValue, err := rconv(src, expectedType.Elem())
		casted := castedValue.Interface()
		if err == nil {
			castedPtrVal := reflect.New(reflect.TypeOf(casted))
			castedPtrVal.Elem().Set(castedValue)
			return castedPtrVal, nil
		} else {
			return src, err
		}
	default:
		return src, nil
	}
}

func rconvAndSet(dvalue reflect.Value, svalue reflect.Value) error {
	tleft := dvalue.Type()
	tright := svalue.Type()
	if tright.AssignableTo(tleft) {
		dvalue.Set(svalue)
	} else {
		vright, err := rconv(svalue, tleft)
		if err == nil {
			dvalue.Set(vright)
		} else {
			return err
		}
	}
	return nil
}

func rconvNumber(value reflect.Value) reflect.Value {
	if !value.IsValid() || value.IsZero() {
		return reflect.ValueOf(0)
	}
	if value.Kind() == reflect.Interface || value.Kind() == reflect.Ptr {
		return rconvNumber(value.Elem())
	}
	if isString(value.Type()) {
		val := value.Interface().(string)
		var num interface{}
		var err error
		if strings.Index(val, ".") > -1 {
			num, err = strconv.ParseFloat(val, 64)
			if err != nil {
				return value
			} else {
				return reflect.ValueOf(num)
			}
		} else {
			num, err = strconv.ParseUint(val, 0, 64)
			num, err = strconv.ParseInt(val, 0, 64)
			if err != nil {
				return value
			} else {
				return reflect.ValueOf(num)
			}
		}
	} else if isBoolean(value.Type()) {
		if value.Bool() {
			return reflect.ValueOf(1)
		} else {
			return reflect.ValueOf(0)
		}
	} else {
		return value
	}
}

func rconvConst(val constant.Value, kind constant.Kind) constant.Value {
	v := constToInterface(val)
	switch kind {
	case constant.Bool:
		return constant.MakeBool(cast.ToBool(v))
	case constant.String:
		return constant.MakeString(cast.ToString(v))
	case constant.Int:
		return constant.MakeInt64(cast.ToInt64(v))
	case constant.Float:
		return constant.MakeFloat64(cast.ToFloat64(v))
	}
	return nil
}

func constToInterface(value constant.Value) interface{} {
	switch value.Kind() {
	case constant.Bool:
		return constant.BoolVal(value)
	case constant.String:
		return constant.StringVal(value)
	case constant.Int:
		v, _ := constant.Int64Val(value)
		return v
	case constant.Float:
		v, _ := constant.Float64Val(value)
		return v
	default:
		return nil
	}
}

func rconvToConcrete(value reflect.Value) reflect.Value {
	if value.Kind() == reflect.Interface {
		return rconvToConcrete(value.Elem())
	} else {
		return value
	}
}

func rconvConstNumber(val constant.Value) (c constant.Value) {
	v := constToInterface(val)
	switch reflect.ValueOf(v).Kind() {
	case reflect.Bool:
		c = constant.MakeInt64(cast.ToInt64(v))
	case reflect.String:
		vstr := v.(string)
		if strings.Index(vstr, ".") > -1 {
			var num float64
			var err error
			num, err = cast.ToFloat64E(vstr)
			if err != nil {
				panic(err)
			}
			c = constant.MakeFloat64(num)
		} else {
			var num int64
			var err error
			num, err = cast.ToInt64E(vstr)
			if err != nil {
				panic(err)
			}
			c = constant.MakeInt64(num)
		}
	case reflect.Int64:
		c = constant.MakeInt64(cast.ToInt64(v))
	case reflect.Float64:
		c = constant.MakeFloat64(cast.ToFloat64(v))
	}
	return
}

func rconvConstInt(value constant.Value) constant.Value {
	v := rconvConstNumber(value)
	if v.Kind() == constant.Float {
		return constant.MakeInt64(cast.ToInt64(constToInterface(v)))
	}
	return v
}

func rconvConstBool(value constant.Value) constant.Value {
	return constant.MakeBool(cast.ToBool(constToInterface(value)))
}

func rconvConstString(value constant.Value) constant.Value {
	return constant.MakeString(cast.ToString(constToInterface(value)))
}

func rconvToString(val reflect.Value) string {
	if !val.IsValid() {
		return ""
	}
	return cast.ToString(val.Interface())
}

func rconvToBool(val reflect.Value) bool {
	if !val.IsValid() {
		return false
	}
	return cast.ToBool(val.Interface())
}

func rconvToInt(val reflect.Value) int {
	if !val.IsValid() {
		return 0
	}
	return cast.ToInt(val.Interface())
}

func rconvToUint(val reflect.Value) uint {
	if !val.IsValid() {
		return 0
	}
	return cast.ToUint(val.Interface())
}

func rconvToInt64(val reflect.Value) int64 {
	if !val.IsValid() {
		return 0
	}
	return cast.ToInt64(val.Interface())
}

func rconvToUint64(val reflect.Value) uint64 {
	if !val.IsValid() {
		return 0
	}
	return cast.ToUint64(val.Interface())
}

func rconvToFloat32(val reflect.Value) float32 {
	if !val.IsValid() {
		return 0
	}
	return cast.ToFloat32(val.Interface())
}

func rconvToFloat64(val reflect.Value) float64 {
	if !val.IsValid() {
		return 0
	}
	return cast.ToFloat64(val.Interface())
}

func rconvToNil(v reflect.Value) bool {
	if isNullable(v.Type()) {
		return v.IsNil() || v.Kind() == reflect.Interface && v.Elem().IsZero()
	} else {
		if v.IsZero() {
			return true
		} else {
			return false
		}
	}
}

func compare(val0, val1 interface{}, op string) (value bool, err error) {
	return rcompare(reflect.ValueOf(val0), reflect.ValueOf(val1), op)
}

func rcompare(val0, val1 reflect.Value, op string) (value bool, err error) {
	if !val0.IsValid() || !val1.IsValid() {
		if !val0.IsValid() && !val1.IsValid() && op == "==" {
			return true, nil
		} else if (val0.IsValid() || val1.IsValid()) && op == "!=" {
			return true, nil
		} else {
			return false, nil
		}
	} else if val0.Kind() == reflect.String || (val0.Kind() == reflect.Interface || val0.Kind() == reflect.Ptr) && val0.Elem().Kind() == reflect.String {
		value, err = compareString(rconvToString(val0), rconvToString(val1), op)
	} else if val0.Kind() == reflect.Bool || (val0.Kind() == reflect.Interface || val0.Kind() == reflect.Ptr) && val0.Elem().Kind() == reflect.Bool {
		value, err = compareBool(rconvToBool(val0), rconvToBool(val1), op)
	} else {
		val0 = rconvNumber(val0)
		val1 = rconvNumber(val1)
		typ0 := val0.Type()
		typ1 := val1.Type()
		if isNumber(typ0) && isNumber(typ1) {
			switch {
			case isUint(typ0):
				switch {
				case isUint(typ1):
					value, err = compareUint(val0.Uint(), val1.Uint(), op)
				case isInt(typ1):
					value, err = compareUint(val0.Uint(), uint64(val1.Int()), op)
				case isFloat(typ1):
					value, err = compareUint(val0.Uint(), uint64(val1.Float()), op)
				}
			case isInt(typ0):
				switch {
				case isUint(typ1):
					value, err = compareInt(val0.Int(), int64(val1.Uint()), op)
				case isInt(typ1):
					value, err = compareInt(val0.Int(), val1.Int(), op)
				case isFloat(typ1):
					value, err = compareInt(val0.Int(), int64(val1.Float()), op)
				}
			case isFloat(typ0):
				switch {
				case isUint(typ1):
					value, err = compareFloat(val0.Float(), float64(val1.Uint()), op)
				case isInt(typ1):
					value, err = compareFloat(val0.Float(), float64(val1.Int()), op)
				case isFloat(typ1):
					value, err = compareFloat(val0.Float(), float64(val1.Float()), op)
				}
			}
		} else {
			err = fmt.Errorf("type %s doesn't support %s operator", typ0.String(), op)
		}
	}
	return
}

func compareString(v0 string, v1 string, op string) (bool, error) {
	switch op {
	case "==":
		return v0 == v1, nil
	case ">":
		return v0 > v1, nil
	case ">=":
		return v0 >= v1, nil
	case "<":
		return v0 < v1, nil
	case "<=":
		return v0 <= v1, nil
	case "!=":
		return v0 != v1, nil
	default:
		return false, fmt.Errorf("unsupported comparison operator %s for string", op)
	}
}

func compareBool(v0 bool, v1 bool, op string) (bool, error) {
	switch op {
	case "==":
		return v0 == v1, nil
	case "!=":
		return v0 != v1, nil
	default:
		return false, fmt.Errorf("unsupported comparison operator %s for bool", op)
	}
}

func compareInt(v0 int64, v1 int64, op string) (bool, error) {
	switch op {
	case "==":
		return v0 == v1, nil
	case ">":
		return v0 > v1, nil
	case ">=":
		return v0 >= v1, nil
	case "<":
		return v0 < v1, nil
	case "<=":
		return v0 <= v1, nil
	case "!=":
		return v0 != v1, nil
	default:
		return false, fmt.Errorf("unsupported comparison operator %s for integer", op)
	}
}

func compareUint(v0 uint64, v1 uint64, op string) (bool, error) {
	switch op {
	case "==":
		return v0 == v1, nil
	case ">":
		return v0 > v1, nil
	case ">=":
		return v0 >= v1, nil
	case "<":
		return v0 < v1, nil
	case "<=":
		return v0 <= v1, nil
	case "!=":
		return v0 != v1, nil
	default:
		return false, fmt.Errorf("unsupported comparison operator %s for unsigned integer", op)
	}
}

func compareFloat(v0 float64, v1 float64, op string) (bool, error) {
	switch op {
	case "==":
		return v0 == v1, nil
	case ">":
		return v0 > v1, nil
	case ">=":
		return v0 >= v1, nil
	case "<":
		return v0 < v1, nil
	case "<=":
		return v0 <= v1, nil
	case "!=":
		return v0 != v1, nil
	default:
		return false, fmt.Errorf("unknown comparison operator %s for float", op)
	}
}

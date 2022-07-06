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

func rconv(val reflect.Value, expected reflect.Type) (reflect.Value, error) {
	valt := val.Type()
	if valt == expected {
		return val, nil
	}
	if valt.Kind() == expected.Kind() && valt.Kind() != reflect.Struct &&
		(valt.PkgPath() != expected.PkgPath() || valt.Name() != expected.Name()) {
		// type def from existing type
		return val.Convert(expected), nil
	}
	value := val.Interface()
	switch expected.Kind() {
	case reflect.Bool:
		casted, err := cast.ToBoolE(value)
		if err == nil {
			return reflect.ValueOf(casted), nil
		} else {
			return val, err
		}
	case reflect.Int:
		casted, err := cast.ToIntE(value)
		if err == nil {
			return reflect.ValueOf(casted), nil
		} else {
			return val, err
		}
	case reflect.Int8:
		casted, err := cast.ToInt8E(value)
		if err == nil {
			return reflect.ValueOf(casted), nil
		} else {
			return val, err
		}
	case reflect.Int16:
		casted, err := cast.ToInt16E(value)
		if err == nil {
			return reflect.ValueOf(casted), nil
		} else {
			return val, err
		}
	case reflect.Int32:
		casted, err := cast.ToInt32E(value)
		if err == nil {
			return reflect.ValueOf(casted), nil
		} else {
			return val, err
		}
	case reflect.Int64:
		casted, err := cast.ToInt64E(value)
		if err == nil {
			return reflect.ValueOf(casted), nil
		} else {
			return val, err
		}
	case reflect.Uint:
		casted, err := cast.ToUintE(value)
		if err == nil {
			return reflect.ValueOf(casted), nil
		} else {
			return val, err
		}
	case reflect.Uint8:
		casted, err := cast.ToUint8E(value)
		if err == nil {
			return reflect.ValueOf(casted), nil
		} else {
			return val, err
		}
	case reflect.Uint16:
		casted, err := cast.ToUint16E(value)
		if err == nil {
			return reflect.ValueOf(casted), nil
		} else {
			return val, err
		}
	case reflect.Uint32:
		casted, err := cast.ToUint32E(value)
		if err == nil {
			return reflect.ValueOf(casted), nil
		} else {
			return val, err
		}
	case reflect.Uint64:
		casted, err := cast.ToUint64E(value)
		if err == nil {
			return reflect.ValueOf(casted), nil
		} else {
			return val, err
		}
	case reflect.Float32:
		casted, err := cast.ToFloat32E(value)
		if err == nil {
			return reflect.ValueOf(casted), nil
		} else {
			return val, err
		}
	case reflect.Float64:
		casted, err := cast.ToFloat64E(value)
		if err == nil {
			return reflect.ValueOf(casted), nil
		} else {
			return val, err
		}
	case reflect.String:
		switch reflect.ValueOf(value).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			casted := fmt.Sprintf("%d", value)
			return reflect.ValueOf(casted), nil
		default:
			casted, err := cast.ToStringE(value)
			if err == nil {
				return reflect.ValueOf(casted), nil
			} else {
				return val, err
			}
		}
	case reflect.Struct:
		castedPtrValue := reflect.New(expected)
		indirect := reflect.Indirect(val)
		kind := indirect.Kind()
		switch kind {
		case reflect.String:
			// assume value is in json format
			err := json.Unmarshal([]byte(indirect.String()), castedPtrValue.Interface())
			if err == nil {
				return castedPtrValue.Elem(), nil
			} else {
				return val, err
			}
		case reflect.Struct:
			err := copier.Copy(castedPtrValue.Interface(), value)
			if err == nil {
				return castedPtrValue.Elem(), nil
			} else {
				return val, err
			}
		case reflect.Interface:
			return indirect.Elem().Convert(expected), nil
		default:
			return val, errors.New(fmt.Sprintf(""))
		}
	case reflect.Map:
		if valt.Kind() != reflect.Map {
			return val, nil
		}
		ktype := expected.Key()
		vtype := expected.Elem()
		castedValue := reflect.MakeMapWithSize(expected, 0)
		keys := val.MapKeys()
		for i := 0; i < len(keys); i++ {
			k := keys[i]
			v := val.MapIndex(k)
			var kcasted, vcasted reflect.Value
			var err error
			kcasted, err = rconv(k, ktype)
			if err != nil {
				return val, err
			}
			vcasted, err = rconv(v, vtype)
			if err != nil {
				return val, err
			}
			castedValue.SetMapIndex(kcasted, vcasted)
		}
		return castedValue, nil
	case reflect.Slice:
		if valt.Kind() != reflect.Slice {
			return val, nil
		}
		vtype := expected.Elem()
		castedValue := reflect.MakeSlice(expected, val.Len(), val.Cap())
		for i := 0; i < val.Len(); i++ {
			vcasted, err := rconv(val.Index(i), vtype)
			if err == nil {
				castedValue.Index(i).Set(vcasted)
			} else {
				return val, err
			}
		}
		return castedValue, nil
	case reflect.Ptr:
		castedValue, err := rconv(val, expected.Elem())
		casted := castedValue.Interface()
		if err == nil {
			castedPtrVal := reflect.New(reflect.TypeOf(casted))
			castedPtrVal.Elem().Set(castedValue)
			return castedPtrVal, nil
		} else {
			return val, err
		}
	default:
		return val, nil
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

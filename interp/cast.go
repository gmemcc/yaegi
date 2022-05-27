package interp

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/spf13/cast"
	"reflect"
)

func maycast(typ *itype, expected *itype) bool {
	if typ.assignableTo(expected) {
		return true
	}
	if typ.rtype.Kind() == reflect.String {
		return true
	}
	ertype := expected.rtype
	_, err := trycast(reflect.New(typ.rtype).Elem(), ertype)
	return err == nil
}

func trycast(val reflect.Value, expected reflect.Type) (reflect.Value, error) {
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
		casted, err := cast.ToStringE(value)
		if err == nil {
			return reflect.ValueOf(casted), nil
		} else {
			return val, err
		}
	case reflect.Struct:
		indirect := reflect.Indirect(val)
		kind := indirect.Kind()
		castedPtrValue := reflect.New(expected)
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
		default:
			return val, errors.New(fmt.Sprintf(""))
		}
	case reflect.Map:
		if val.Type().Kind() != reflect.Map {
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
			kcasted, err = trycast(k, ktype)
			if err != nil {
				return val, err
			}
			vcasted, err = trycast(v, vtype)
			if err != nil {
				return val, err
			}
			castedValue.SetMapIndex(kcasted, vcasted)
		}
		return castedValue, nil
	case reflect.Slice:
		if val.Type().Kind() != reflect.Slice {
			return val, nil
		}
		vtype := expected.Elem()
		castedValue := reflect.MakeSlice(expected, val.Len(), val.Cap())
		for i := 0; i < val.Len(); i++ {
			vcasted, err := trycast(val.Index(i), vtype)
			if err == nil {
				castedValue.Index(i).Set(vcasted)
			} else {
				return val, err
			}
		}
		return castedValue, nil
	case reflect.Ptr:
		castedValue, err := trycast(val, expected.Elem())
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

package interp

import (
	"fmt"
	"reflect"
	"testing"
)

var (
	typeofInt     = reflect.TypeOf(0)
	typeofInt8    = reflect.TypeOf(int8(0))
	typeofInt16   = reflect.TypeOf(int16(0))
	typeofInt32   = reflect.TypeOf(int32(0))
	typeofInt64   = reflect.TypeOf(int64(0))
	typeofUint    = reflect.TypeOf(uint(0))
	typeofUint8   = reflect.TypeOf(uint8(0))
	typeofUint16  = reflect.TypeOf(uint16(0))
	typeofUint32  = reflect.TypeOf(uint32(0))
	typeofUint64  = reflect.TypeOf(uint64(0))
	typeofFloat32 = reflect.TypeOf(float32(0))
	typeofFloat64 = reflect.TypeOf(float64(0))
	typeofString  = reflect.TypeOf("")
	typeofBool    = reflect.TypeOf(false)
	typeofAny     = reflect.TypeOf((*interface{})(nil)).Elem()
	typeofError   = reflect.TypeOf((*error)(nil)).Elem()
)

func TestRconv(t *testing.T) {
	suint8 := []uint8{65, 66, 67, 68, 69}
	dorconv(reflect.ValueOf(suint8), typeofString)
	sfloat64 := []float64{1.2, 2.3, 3.4}
	dorconv(reflect.ValueOf(sfloat64), reflect.TypeOf([]int{}))

	bfalse := false
	i100 := 100
	s200 := "200"

	dorconv(reflect.ValueOf(bfalse), typeofString)
	dorconv(reflect.ValueOf(i100), typeofBool)
	dorconv(reflect.ValueOf(s200), typeofInt)

	var iany interface{}
	iany = bfalse
	dorconv(reflect.ValueOf(iany), typeofString)
	iany = i100
	dorconv(reflect.ValueOf(iany), typeofBool)
	iany = s200
	dorconv(reflect.ValueOf(iany), typeofInt)
}

func dorconv(src reflect.Value, expectedType reflect.Type) {
	result, err := rconv(src, expectedType)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v %#v => %v %#v\n", src.Type().String(), src, result.Type().String(), result)
}

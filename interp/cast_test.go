package interp

import (
	"fmt"
	"reflect"
	"testing"
)

type user struct {
	Name string `json:"name"`
}

type person struct {
	Name string `json:"name"`
}

func Test_tryCast(t *testing.T) {
	type args struct {
		val      reflect.Value
		expected reflect.Type
	}
	intNum := 123
	floatNum := float32(123)
	str := "123"
	p := person{Name: "alex"}
	strslice := []string{"1", "2", "3"}
	tests := []struct {
		name string
		args args
		want reflect.Value
	}{
		{
			name: "primitive int",
			args: args{
				val:      reflect.ValueOf(intNum),
				expected: reflect.TypeOf(""),
			},
			want: reflect.ValueOf(str),
		},
		{
			name: "primitive float",
			args: args{
				val:      reflect.ValueOf(floatNum),
				expected: reflect.TypeOf(""),
			},
			want: reflect.ValueOf(str),
		},
		{
			name: "pointer",
			args: args{
				val:      reflect.ValueOf(intNum),
				expected: reflect.TypeOf(&str),
			},
			want: reflect.ValueOf(&str),
		},
		{
			name: "struct",
			args: args{
				val:      reflect.ValueOf(user{Name: "alex"}),
				expected: reflect.TypeOf(p),
			},
			want: reflect.ValueOf(p),
		},
		{
			name: "str2struct",
			args: args{
				val:      reflect.ValueOf(`{ "name": "alex" }`),
				expected: reflect.TypeOf(p),
			},
			want: reflect.ValueOf(p),
		},
		{
			name: "slice",
			args: args{
				val:      reflect.ValueOf([]interface{}{"1", "2", "3"}),
				expected: reflect.TypeOf(strslice),
			},
			want: reflect.ValueOf(strslice),
		},
		{
			name: "slice",
			args: args{
				val:      reflect.ValueOf([]interface{}{"1", "2", "3"}),
				expected: reflect.TypeOf(strslice),
			},
			want: reflect.ValueOf(strslice),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf("value to cast %#v expected type %#v\n", tt.args.val.Interface(), tt.args.expected.String())
			got, err := trycast(tt.args.val, tt.args.expected)
			if err != nil {
				t.Errorf("%s", err.Error())
			}
			goti := reflect.Indirect(got).Interface()
			wanti := reflect.Indirect(tt.want).Interface()
			gott := got.Type().String()
			wantt := tt.want.Type().String()
			if !reflect.DeepEqual(goti, wanti) || gott != wantt {
				t.Errorf("trycast() got = %#v, want %#v. got type = %s, want type = %s", got, tt.want, gott, wantt)
			}
		})
	}
}

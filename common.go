package smartpb

import "reflect"

func indirect(v reflect.Value) reflect.Value {
	if v.Kind() != reflect.Pointer {
		return v
	}
	return indirect(v.Elem())
}

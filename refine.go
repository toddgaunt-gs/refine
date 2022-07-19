package refine

import (
	"reflect"
)

const tag = "refine"

func Check(val any) {
	t := func() reflect.Type {
		t := reflect.TypeOf(val)
		if t.Kind() == reflect.Pointer {
			return t.Elem()
		} else {
			return t
		}
	}()

	n := t.NumField()
	for i := 0; i < n; i++ {
		field := t.Field(i)
		var _ = field.Tag.Get(tag)
	}
}

func eval(expr *expression) {
}

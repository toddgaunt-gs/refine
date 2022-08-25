package refine

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrNotStruct       = errors.New("only struct types can be checked")
	ErrUnsupportedType = errors.New("unsupported type")
	ErrNotMet          = errors.New("not met")
	ErrParse           = errors.New("could not be parsed")
	ErrEval            = errors.New("could not be evaluated")
)

type checkErr struct {
	structType string
	fieldName  string
	fieldValue any
	refinement string
	err        error
}

func (e checkErr) Error() string {
	return fmt.Sprintf(
		"refine.Check: %s.%s = %#+v, %q %v",
		e.structType, e.fieldName, e.fieldValue, e.refinement, e.err,
	)
}

func (e checkErr) Unwrap() error {
	return e.err
}

const tag = "refine"

var kindMap = map[reflect.Kind]Kind{
	reflect.String:  boxString,
	reflect.Bool:    boxBool,
	reflect.Int:     boxInt,
	reflect.Int8:    boxInt,
	reflect.Int16:   boxInt,
	reflect.Int32:   boxInt,
	reflect.Int64:   boxInt,
	reflect.Uint:    boxUint,
	reflect.Uint8:   boxUint,
	reflect.Uint16:  boxUint,
	reflect.Uint32:  boxUint,
	reflect.Uint64:  boxUint,
	reflect.Slice:   boxSlice,
	reflect.Map:     boxMap,
	reflect.Pointer: boxPointer,
}

func Check(val any) error {
	t := func() reflect.Type {
		t := reflect.TypeOf(val)
		if t.Kind() == reflect.Pointer {
			return t.Elem()
		} else {
			return t
		}
	}()

	v := func() reflect.Value {
		v := reflect.ValueOf(val)
		if v.Kind() == reflect.Pointer {
			return v.Elem()
		} else {
			return v
		}
	}()

	if t.Kind() != reflect.Struct || v.Kind() != reflect.Struct {
		return fmt.Errorf("refine.Check: %s is a %s, %w", t.Name(), v.Kind().String(), ErrNotStruct)
	}

	var ev = newEvaluator()

	n := t.NumField()

	// Populate the symbol table with values from the struct's fields.
	for i := 0; i < n; i++ {
		field := t.Field(i)
		kind := field.Type.Kind()
		value := v.Field(i)

		boxKind, ok := kindMap[kind]
		if !ok {
			return fmt.Errorf("refine.Check: %s.%s %w %s", t.Name(), field.Name, ErrUnsupportedType, kind.String())
		}

		if kind == reflect.Pointer {
			if value.IsNil() {
				ev.symbols[field.Name] = Value{
					kind: boxPointer,
					val:  nil,
				}
			} else {
				ev.symbols[field.Name] = Value{
					kind: boxKind,
					val:  value.Interface(),
				}
			}
		} else {
			ev.symbols[field.Name] = Value{
				kind: boxKind,
				val:  value.Interface(),
			}
		}
	}

	// Parse and evaluate the refinements on each field.
	for i := 0; i < n; i++ {
		field := t.Field(i)
		refinement := field.Tag.Get(tag)
		tokens := lex(field.Name, refinement)

		expr, err := parse(tokens)
		if err != nil {
			return checkErr{
				structType: t.Name(),
				fieldName:  field.Name,
				fieldValue: ev.symbols[field.Name].val,
				refinement: refinement,
				err:        fmt.Errorf("%w: %v", ErrParse, err),
			}
		}

		//result, err := eval(ev, expr)
		expr.Accept(ev)
		result, err := ev.Result, ev.Err
		if err != nil {
			return checkErr{
				structType: t.Name(),
				fieldName:  field.Name,
				fieldValue: ev.symbols[field.Name].val,
				refinement: refinement,
				err:        fmt.Errorf("%w: %v", ErrEval, err),
			}
		}

		if result.kind != boxBool {
			return checkErr{
				structType: t.Name(),
				fieldName:  field.Name,
				fieldValue: ev.symbols[field.Name].val,
				refinement: refinement,
				err:        fmt.Errorf("%w: %v", ErrEval, "not bool"),
			}
		}

		if result.val.(bool) != true {
			return checkErr{
				structType: t.Name(),
				fieldName:  field.Name,
				fieldValue: ev.symbols[field.Name].val,
				refinement: refinement,
				err:        ErrNotMet,
			}
		}
	}

	return nil
}

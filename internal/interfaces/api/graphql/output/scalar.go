package output

import (
	"reflect"
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

var Void = graphql.NewScalar(graphql.ScalarConfig{
	Name:         "void",
	Serialize:    func(value interface{}) interface{} { return nil },
	ParseValue:   func(value interface{}) interface{} { return nil },
	ParseLiteral: func(valueAST ast.Value) interface{} { return nil },
})

var AnyInt = graphql.NewScalar(graphql.ScalarConfig{
	Name:       "AnyInt",
	Serialize:  coerceAnyInt,
	ParseValue: coerceAnyInt,
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.IntValue:
			if intVal, err := strconv.ParseInt(valueAST.Value, 10, 64); err == nil {
				return intVal
			}
		}
		return nil
	},
})

func coerceAnyInt(val interface{}) interface{} {
	v := reflect.ValueOf(val)
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
		val = v.Interface()
	}

	switch v := val.(type) {
	case bool:
		if v {
			return 1
		}
		return 0
	case int, int8, int16, int32, int64:
		return int(reflect.ValueOf(v).Int())
	case uint, uint8, uint16, uint32, uint64:
		return int(reflect.ValueOf(v).Uint())
	case float32, float64:
		return int(reflect.ValueOf(v).Float())
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return int(f)
		}
	}
	return nil
}

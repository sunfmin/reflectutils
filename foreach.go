package reflectutils

import (
	"fmt"
	"reflect"
)

func ForEach(arr interface{}, predicate interface{}) {
	var (
		funcValue = reflect.ValueOf(predicate)
		arrValue  = reflect.ValueOf(arr)
		arrType   = arrValue.Type()
		arrKind   = arrType.Kind()
		funcType  = funcValue.Type()
	)

	if IsNil(arr) {
		return
	}

	if !(arrKind == reflect.Array || arrKind == reflect.Slice || arrKind == reflect.Map) {
		panic("First parameter must be an iteratee")
	}

	if arrKind == reflect.Slice || arrKind == reflect.Array {
		if !isFunction(predicate, 1, 0) {
			panic("Second argument must be a function with one parameter")
		}

		arrElemType := arrValue.Type().Elem()

		// Checking whether element type is convertible to function's first argument's type.
		if !arrElemType.ConvertibleTo(funcType.In(0)) {
			panic("Map function's argument is not compatible with type of array.")
		}

		for i := 0; i < arrValue.Len(); i++ {
			funcValue.Call([]reflect.Value{arrValue.Index(i)})
		}
		return
	}

	if arrKind == reflect.Map {
		if !isFunction(predicate, 2, 0) {
			panic("Second argument must be a function with two parameters")
		}

		// Type checking for Map<key, value> = (key, value)
		keyType := arrType.Key()
		valueType := arrType.Elem()

		if !keyType.ConvertibleTo(funcType.In(0)) {
			panic(fmt.Sprintf("function first argument is not compatible with %s", keyType.String()))
		}

		if !valueType.ConvertibleTo(funcType.In(1)) {
			panic(fmt.Sprintf("function second argument is not compatible with %s", valueType.String()))
		}

		for _, key := range arrValue.MapKeys() {
			funcValue.Call([]reflect.Value{key, arrValue.MapIndex(key)})
		}
		return
	}
}

// isFunction returns if the argument is a function.
func isFunction(in interface{}, num ...int) bool {
	funcType := reflect.TypeOf(in)

	result := funcType != nil && funcType.Kind() == reflect.Func

	if len(num) >= 1 {
		result = result && funcType.NumIn() == num[0]
	}

	if len(num) == 2 {
		result = result && funcType.NumOut() == num[1]
	}

	return result
}

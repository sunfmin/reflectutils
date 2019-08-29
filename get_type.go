package reflectutils

import (
	"reflect"
	"strings"
)

// Get value of a struct by path using reflect.
func GetType(i interface{}, name string) (t reflect.Type) {
	var err error

	t = reflect.TypeOf(i)

	if name == "" {
		return
	}

	var token *dotToken
	token, err = nextDot(name)
	if err != nil {
		return nil
	}

	if t.Kind() == reflect.Map || t.Kind() == reflect.Slice {
		t = GetType(reflect.Zero(t.Elem()).Interface(), token.Left)
		return
	}

	if t.Kind() != reflect.Struct {
		for t.Elem().Kind() == reflect.Ptr {
			t = t.Elem()
		}

		t = t.Elem()
	}

	if t.Kind() == reflect.Struct {

		sf, ok := t.FieldByNameFunc(func(name string) bool {
			return strings.EqualFold(name, token.Field)
		})

		if !ok {
			return nil
		}

		t = GetType(reflect.Zero(sf.Type).Interface(), token.Left)
		return
	}

	return
}

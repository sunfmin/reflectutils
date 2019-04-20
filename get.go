package reflectutils

import (
	"fmt"
	"reflect"
	"strings"
)

// MustGet get value of a struct by path using reflect, return nil if anything in the path is nil
func MustGet(i interface{}, name string) (value interface{}) {
	var err error
	value, err = Get(i, name)
	if err != nil {
		panic(err)
	}
	return
}

// Get value of a struct by path using reflect.
func Get(i interface{}, name string) (value interface{}, err error) {
	// printv(i, name)

	if name == "" {
		value = i
		return
	}

	v := reflect.ValueOf(i)

	var token *dotToken
	token, err = nextDot(name)
	if err != nil {
		return
	}

	sv := v

	if sv.Kind() == reflect.Map {
		// map must have string type
		mv := sv

		if mv.Type().Key() != reflect.TypeOf("") {
			err = fmt.Errorf("map key %s must be string type", name)
			return
		}

		if mv.IsNil() {
			return
		}

		keyValue := reflect.ValueOf(token.Field)

		elemType := mv.Type().Elem()
		mapElem := reflect.New(elemType).Elem()
		existElem := mv.MapIndex(keyValue)
		if existElem.IsValid() {
			mapElem.Set(existElem)
		}

		value, err = Get(mapElem.Interface(), token.Left)
		if err != nil {
			return
		}
		return
	}

	if sv.Kind() == reflect.Slice {
		av := sv

		if token.IsAppendingArray {
			err = fmt.Errorf("array index is empty: %s", name)
			return
		}

		if av.Len() <= token.ArrayIndex {
			return
		}

		arrayElem := av.Index(token.ArrayIndex)
		if !arrayElem.IsValid() {
			return
		}

		value, err = Get(arrayElem.Interface(), token.Left)
		if err != nil {
			return
		}

		return
	}

	if sv.Kind() != reflect.Struct {
		for sv.Elem().Kind() == reflect.Ptr {
			sv = sv.Elem()
		}

		sv = sv.Elem()
	}

	if sv.Kind() == reflect.Struct {
		fv := sv.FieldByNameFunc(func(fname string) bool {
			return strings.EqualFold(fname, token.Field)
		})

		if !fv.IsValid() {
			err = NoSuchFieldError
			return
		}
		value, err = Get(fv.Interface(), token.Left)
		return
	}

	return
}

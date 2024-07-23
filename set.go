/*
# reflectutils

reflectutils is for setting your struct value by using string name and path that follows reflectutils rules.

## The rules

- `.Name` to set a property by field name
- `.Person.Name` to set the name of the current property
- `.Person.Addresses[0].Phone` to set an element of an array property
- `.Person.Addresses[].Name` it will create a object of address and set it's property
- `.Person.MapData.Name` it can also set value to map

## How to install

	go get github.com/sunfmin/reflectutils
*/
package reflectutils

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

var NoSuchFieldError = errors.New("no such field")

// Set value of a struct by path using reflect.
func Set(i interface{}, name string, value interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
		}
	}()

	v := reflect.ValueOf(i)

	if v.Kind() != reflect.Ptr {
		err = errors.New("set object must be a pointer.")
		return
	}

	for v.Elem().Kind() == reflect.Ptr {
		v = v.Elem()
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
	}

	sv := v.Elem()

	if name == "" {

		switch inputv := value.(type) {
		case string:
			err = setStringValue(sv, inputv)
			if err != nil {
				return
			}
		case []byte:
			err = setStringValue(sv, string(inputv))
			if err != nil {
				return
			}
		default:
			if value == nil {
				vm := reflect.ValueOf(i)
				vm.Elem().Set(reflect.Zero(vm.Elem().Type()))
				return
			}

			valv := reflect.ValueOf(value)
			for valv.Kind() == reflect.Ptr {
				valv = valv.Elem()
			}
			sv.Set(valv)
		}

		return
	}

	var token *dotToken
	token, err = nextDot(name)
	if err != nil {
		return
	}

	// printv(sv.Interface(), name, value)

	if sv.Kind() == reflect.Map {
		// map must have string type
		mv := sv

		if mv.Type().Key() != reflect.TypeOf("") {
			return fmt.Errorf("map key %s must be string type", name)
		}

		if mv.IsNil() {
			mv.Set(reflect.MakeMap(mv.Type()))
		}

		keyValue := reflect.ValueOf(token.Field)

		elemType := mv.Type().Elem()
		mapElem := reflect.New(elemType).Elem()

		existElem := mv.MapIndex(keyValue)
		if existElem.IsValid() {
			mapElem.Set(existElem)
		}

		err = Set(mapElem.Addr().Interface(), token.Left, value)
		if err != nil {
			return
		}

		mv.SetMapIndex(keyValue, mapElem)
		return
	}

	if sv.Kind() == reflect.Slice {
		av := sv
		elemType := av.Type().Elem()
		var newslice reflect.Value

		if token.IsAppendingArray {
			newslice = av
			arrayElem := reflect.New(elemType).Elem()
			err = Set(arrayElem.Addr().Interface(), token.Left, value)
			if err != nil {
				return
			}
			newslice = reflect.Append(newslice, arrayElem)
		} else {
			if av.Len() > token.ArrayIndex {
				newslice = reflect.MakeSlice(av.Type(), 0, 0)

				arrayElem := av.Index(token.ArrayIndex)
				if !arrayElem.IsValid() {
					arrayElem.Set(reflect.New(elemType).Elem())
				}

				err = Set(arrayElem.Addr().Interface(), token.Left, value)
				if err != nil {
					return
				}
				for i := 0; i < token.ArrayIndex; i++ {
					newslice = reflect.Append(newslice, av.Index(i))
				}
				newslice = reflect.Append(newslice, arrayElem)
				for i := token.ArrayIndex + 1; i < av.Len(); i++ {
					newslice = reflect.Append(newslice, av.Index(i))
				}
			} else {
				newslice = av
				if newslice.IsNil() {
					newslice = reflect.MakeSlice(newslice.Type(), 0, 0)
				}
				arrayElem := reflect.New(elemType).Elem()
				err = Set(arrayElem.Addr().Interface(), token.Left, value)
				if err != nil {
					return
				}
				if newslice.Len() < token.ArrayIndex {
					for newslice.Len() < token.ArrayIndex {
						newslice = reflect.Append(newslice, reflect.Zero(elemType))
					}
				}

				newslice = reflect.Append(newslice, arrayElem)

			}
		}

		av.Set(newslice)

		return
	}

	if sv.Kind() == reflect.Struct {
		fv := sv.FieldByNameFunc(func(fname string) bool {
			return strings.EqualFold(fname, token.Field)
		})

		if !fv.IsValid() {
			// err = errors.New(fmt.Sprintf("%+v has no such field `%s`.", sv.Interface(), token.Field))
			err = NoSuchFieldError
			return
		}

		err = Set(fv.Addr().Interface(), token.Left, value)
		return
	}

	return
}

type dotToken struct {
	Field            string
	Left             string
	IsArray          bool
	ArrayIndex       int
	IsAppendingArray bool
}

func nextDot(name string) (t *dotToken, err error) {
	t = &dotToken{}
	t.Field = strings.Trim(name, ".[")

	if i := strings.IndexAny(t.Field, ".["); i > 0 {
		t.Field, t.Left = t.Field[:i], t.Field[i+1:]
	}

	if t.Field[len(t.Field)-1:] == "]" {
		t.IsArray = true
		arrayIndexString := t.Field[0 : len(t.Field)-1]
		if arrayIndexString == "" {
			t.IsAppendingArray = true
		} else {
			var i64 int64
			i64, err = strconv.ParseInt(arrayIndexString, 10, 64)
			t.ArrayIndex = int(i64)
			if err != nil {
				return
			}
		}
	}

	return
}

func printv(v interface{}, name interface{}) {
	log.Println("=====")
	rv := reflect.ValueOf(v)
	log.Printf(
		"\n\tname: %+v, \n\tv: %+v, \n\trv: %+v, \n\trv.Kind(): %+v, \n\trv.Type(): %+v, \n\trv.IsValid(): %+v",
		name,
		v,
		rv,
		rv.Kind(),
		rv.Type(),
		rv.IsValid(),
	)
	log.Println("=====")
}

func setStringValue(v reflect.Value, value string) (err error) {
	s := value

	// if type is []byte
	if v.Kind() == reflect.Slice && v.Type().Elem().Kind() == reflect.Uint8 {
		v.SetBytes([]byte(s))
		return
	}

	switch v.Kind() {
	case reflect.String:
		v.SetString(s)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var n int64
		n, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			return
		}
		if v.OverflowInt(n) {
			err = fmt.Errorf("overflow int64 for %d", n)
			return
		}
		v.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		var n uint64
		n, err = strconv.ParseUint(s, 10, 64)
		if err != nil {
			return
		}
		if v.OverflowUint(n) {
			err = fmt.Errorf("overflow uint64 for %d", n)
			return
		}
		v.SetUint(n)
	case reflect.Float32, reflect.Float64:
		var n float64
		n, err = strconv.ParseFloat(s, v.Type().Bits())
		if err != nil {
			return
		}
		if v.OverflowFloat(n) {
			err = fmt.Errorf("overflow float64 for %f", n)
			return
		}
		v.SetFloat(n)
	case reflect.Bool:
		var n bool
		n, err = strconv.ParseBool(s)
		if err != nil {
			return
		}
		v.SetBool(n)
	default:
		err = fmt.Errorf("value %+v can only been set to primary type but was %+v", value, v)
	}

	return
}

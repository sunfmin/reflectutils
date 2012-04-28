package reflectutils

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

func Set(i interface{}, name string, value string) (err error) {

	v := reflect.ValueOf(i)

	if v.Kind() != reflect.Ptr {
		panic("set object must be a pointer.")
	}

	for v.Elem().Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}

	sv := v.Elem()

	if name == "" {
		setStringValue(sv, value)
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
			return errors.New(fmt.Sprintf("map key %s must be string type", name))
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
			// err = errors.New(fmt.Sprintf("can not find field %s.", field))
			return errors.New(fmt.Sprintf("can not find field %s.", token.Left))
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

func printv(v interface{}, name interface{}, value string) {
	log.Println("=====")
	rv := reflect.ValueOf(v)
	log.Printf(
		"\n\tname: %+v, \n\tv: %+v, \n\trv: %+v, \n\trv.Kind(): %+v, \n\trv.Type(): %+v, \n\trv.IsNil(): %+v, \n\trv.IsValid(): %+v",
		name,
		v,
		rv,
		rv.Kind(),
		rv.Type(),
		"",
		rv.IsValid(),
	)
	log.Println("=====\n\n")
}

func Populate(i interface{}) {
	v := reflect.ValueOf(i)

	if v.Kind() != reflect.Ptr {
		panic("object must be a pointer.")
	}

	for v.Elem().Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}

	// printv(v.Interface(), "", "")

	sv := v.Elem()

	for i := 0; i < sv.NumField(); i++ {
		f := sv.Field(i)
		if f.Kind() == reflect.Ptr && f.IsNil() {
			Populate(f.Addr().Interface())
		}
		// printv(f.Interface(), st.Field(i).Name, "")
	}

	// printv(i, "", "")
}

func setStringValue(v reflect.Value, value string) (err error) {
	s := value
	switch v.Kind() {
	case reflect.String:
		v.SetString(s)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil || v.OverflowInt(n) {
			break
		}
		v.SetInt(n)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil || v.OverflowUint(n) {
			break
		}
		v.SetUint(n)

	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(s, v.Type().Bits())
		if err != nil || v.OverflowFloat(n) {
			break
		}
		v.SetFloat(n)
	default:
		panic(fmt.Sprintf("must set primary type but was %+v", v))
	}

	return
}

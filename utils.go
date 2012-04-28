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
		fmt.Println("Set primary start ===>", sv.Interface(), value)
		setStringValue(sv, value)
		return
	}

	fmt.Println("Set object start ===>")
	var token *dotToken
	token, err = nextDot(name)
	if err != nil {
		return
	}

	printv(sv.Interface(), name, value)

	if sv.Kind() == reflect.Map {
		// map must have string type
		mv := sv

		if mv.Type().Key() != reflect.TypeOf("") {
			return errors.New(fmt.Sprintf("map key %s must be string type", name))
		}

		if mv.IsNil() {
			mv.Set(reflect.MakeMap(mv.Type()))
		}

		var mapElem reflect.Value
		elemType := mv.Type().Elem()

		if !mapElem.IsValid() {
			mapElem = reflect.New(elemType).Elem()
		} else {
			mapElem.Set(reflect.Zero(elemType))
		}

		err = Set(mapElem.Addr().Interface(), token.Left, value)
		if err != nil {
			return
		}

		mv.SetMapIndex(reflect.ValueOf(token.Field), mapElem)
		return
	}

	if sv.Kind() == reflect.Slice {
		av := sv
		if av.IsNil() {
			av.Set(reflect.MakeSlice(av.Type(), 0, 0))
		}

		var arrayElem reflect.Value
		elemType := av.Type().Elem()

		if !arrayElem.IsValid() {
			arrayElem = reflect.New(elemType).Elem()
		} else {
			arrayElem.Set(reflect.Zero(elemType))
		}

fmt.Println(token.Left)

		err = Set(arrayElem.Addr().Interface(), token.Left, value)
		if err != nil {
			return
		}
		av = reflect.Append(av, arrayElem)
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
	ArrayIndex       uint64
	IsAppendingArray bool
}

func nextDot(name string) (t *dotToken, err error) {
	t = &dotToken{}
	t.Field = strings.Trim(name, ".")
	t.Left = ""

	if i := strings.IndexAny(t.Field, "."); i > 0 {
		t.Field, t.Left = t.Field[:i], t.Field[i+1:]
	}

	if j := strings.IndexAny(t.Field, "["); j > 0 {
		t.IsArray = true
		arrayIndexBracket := ""
		t.Field, arrayIndexBracket = t.Field[:j], t.Field[j+1:]
		if arrayIndexBracket[len(arrayIndexBracket)-1:] != "]" {
			err = errors.New(fmt.Sprintf("missing ] for %v", name))
			return
		}
		arrayIndexString := arrayIndexBracket[0 : len(arrayIndexBracket)-1]
		if arrayIndexString == "" {
			t.IsAppendingArray = true
		} else {
			t.ArrayIndex, err = strconv.ParseUint(arrayIndexString, 10, 64)
			if err != nil {
				return
			}
		}
	}

	fmt.Println("token:", t)
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

	printv(v.Interface(), "", "")

	sv := v.Elem()
	st := sv.Type()

	for i := 0; i < sv.NumField(); i++ {
		f := sv.Field(i)
		if f.Kind() == reflect.Ptr && f.IsNil() {
			Populate(f.Addr().Interface())
		}
		printv(f.Interface(), st.Field(i).Name, "")
	}

	printv(i, "", "")
}

func setStringValue(v reflect.Value, value string) (err error) {
	s := value
	log.Println("setValue:", v.Kind(), value)

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

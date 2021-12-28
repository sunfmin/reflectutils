package reflectutils

import (
	"reflect"
	"strconv"
	"strings"
)

func Delete(i interface{}, name string) (err error) {
	key := ""
	if strings.HasSuffix(name, "]") {
		lb := strings.LastIndex(name, "[")
		key = name[lb+1 : len(name)-1]
		name = name[0:lb]
	}

	t := GetType(i, name)
	if t == nil {
		return NoSuchFieldError
	}

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() == reflect.Slice {
		var index int
		index, err = strconv.Atoi(key)
		if err != nil {
			return
		}
		newSlice := reflect.MakeSlice(t, 0, 0)
		v, err1 := Get(i, name)
		if err1 != nil {
			return err1
		}
		vv := reflect.ValueOf(v)
		for vv.Kind() == reflect.Ptr {
			vv = vv.Elem()
		}
		for j := 0; j < vv.Len(); j++ {
			if j == index {
				continue
			}
			newSlice = reflect.Append(newSlice, vv.Index(j))
		}
		return Set(i, name, newSlice.Interface())
	}

	if t.Kind() == reflect.Map {
		v := MustGet(i, name)
		vv := reflect.ValueOf(v)
		for vv.Kind() == reflect.Ptr {
			vv = vv.Elem()
		}
		vv.SetMapIndex(reflect.ValueOf(key), reflect.Value{})
		return
	}

	if t.Kind() == reflect.Struct {
		return Set(i, name, nil)
	}

	return Set(i, name, reflect.Zero(t).Interface())
}

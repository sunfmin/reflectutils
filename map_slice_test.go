package reflectutils_test

import (
	"testing"

	. "github.com/sunfmin/reflectutils"
)

type mapTest struct {
	Name   string
	Value  string
	Getter func(db map[string][]map[string]string) string
}

var mapcases = []*mapTest{
	{
		Name:  "users[0].gender",
		Value: "male",
		Getter: func(db map[string][]map[string]string) string {
			return db["users"][0]["gender"]
		},
	},
	{
		Name:  "users[1].gender",
		Value: "female",
		Getter: func(db map[string][]map[string]string) string {
			return db["users"][1]["gender"]
		},
	},
	{
		Name:  "companies[].name",
		Value: "The Plant",
		Getter: func(db map[string][]map[string]string) string {
			return db["companies"][0]["name"]
		},
	},
	{
		Name:  "companies[].name",
		Value: "AQ",
		Getter: func(db map[string][]map[string]string) string {
			return db["companies"][1]["name"]
		},
	},
}

func TestSetMapAndSlice(t *testing.T) {
	var db map[string][]map[string]string

	for _, c := range mapcases {
		err := Set(&db, c.Name, c.Value)
		if err != nil {
			t.Error(err)
		}
		val := c.Getter(db)
		if val != c.Value {
			t.Errorf("expected is %v, but was %v", c.Value, val)
		}
	}

}

func TestSlices(t *testing.T) {
	var val []string
	Set(&val, "[100]", "100")
	if val[0] != "" {
		t.Error("val[0] is not empty")
	}

	if val[100] != "100" {
		t.Error("val[100] is not 100")
	}

}

func TestMapSetNil(t *testing.T) {
	m := make(map[string]int)
	Set(&m, "", nil)
	if m != nil {
		t.Errorf("got non-nil (%p), want nil", m)
	}
}

func TestSliceSetNil(t *testing.T) {
	m := []int{1}
	err := Set(&m, "", nil)
	if err != nil {
		panic(err)
	}
	if m != nil {
		t.Errorf("got non-nil (%p), want nil", m)
	}

}

func TestStructSetNil(t *testing.T) {
	type S struct {
		Val string
	}
	m := &S{Val: "123"}
	err := Set(&m, "", nil)
	if err != nil {
		panic(err)
	}


	if m != nil {
		t.Errorf("got non-nil (%#+v), want nil", m)
	}
}

func TestSetFieldNil(t *testing.T) {

	type S struct {
		Value []int
	}

	var s = &S{
		Value: []int{1, 2},
	}

	err := Set(s, "Value", nil)
	if err != nil {
		panic(err)
	}

	if s.Value != nil {
		panic("s.Value is not nil")
	}

}
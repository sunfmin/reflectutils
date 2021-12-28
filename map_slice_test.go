package reflectutils_test

import (
	"fmt"
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

func TestDelete(t *testing.T) {
	p := &Person{
		Departments: []*Department{
			{
				Name: "D1",
			},
			{
				Name: "D2",
			},
			{
				Name: "D3",
			},
		},
		Phones: map[string]string{
			"felix": "123",
			"john":  "456",
		},
		Projects: []*Project{
			{
				Name: "Project1",
				Members: []*Person{
					{
						Name: "Felix1",
					},
					{
						Name: "Felix2",
					},
				},
			},
		},
	}
	var strs = []string{"1", "2", "3"}
	var map1 = map[string]string{"a": "a1", "b": "b1"}
	var map2 = &map1

	var p2 = &Person{}

	var deleteCases = []struct {
		caseName    string
		obj         interface{}
		name        string
		result      func(obj interface{}) (v interface{})
		expected    string
		expectedErr string
	}{
		{
			caseName: "Delete map element in object",
			obj:      p,
			name:     "Phones[felix]",
			result: func(obj interface{}) (v interface{}) {
				return obj.(*Person).Phones
			},
			expected: "map[john:456]",
		},
		{
			caseName: "Delete slice element in object",
			obj:      p,
			name:     "Departments[1]",
			result: func(obj interface{}) (v interface{}) {
				return obj.(*Person).Departments[1]
			},
			expected: "&{Id:0 Name:D3}",
		},

		{
			caseName: "Delete slice element in object nested",
			obj:      p,
			name:     "Projects[0].Members[0]",
			result: func(obj interface{}) (v interface{}) {
				return obj.(*Person).Projects[0].Members[0].Name + obj.(*Person).Projects[0].Name
			},
			expected: "Felix2Project1",
		},

		{
			caseName: "Delete slice element in object overflow",
			obj:      p,
			name:     "Departments[100]",
			result: func(obj interface{}) (v interface{}) {
				return obj.(*Person).Departments[1]
			},
			expected: "&{Id:0 Name:D3}",
		},

		{
			caseName:    "Delete struct with error",
			obj:         p,
			name:        "Comp21",
			expectedErr: "no such field",
		},

		{
			caseName: "Delete struct",
			obj:      p,
			name:     "Departments[0].Name",
			result: func(obj interface{}) (v interface{}) {
				return obj.(*Person).Departments[0].Name
			},
			expected: "",
		},

		{
			caseName:    "Delete with wrong index",
			obj:         p,
			name:        "Departments[abc].Name",
			expectedErr: "no such field",
		},

		{
			caseName: "Delete slice element",
			obj:      &strs,
			name:     "[1]",
			result: func(obj interface{}) (v interface{}) {
				return strs
			},
			expected: `[1 3]`,
		},

		{
			caseName: "Delete map element",
			obj:      &map2,
			name:     "[a]",
			result: func(obj interface{}) (v interface{}) {
				return map1
			},
			expected: `map[b:b1]`,
		},

		{
			caseName: "Delete struct whole",
			obj:      &p2,
			name:     "",
			result: func(obj interface{}) (v interface{}) {
				return p2
			},
			expected: `<nil>`,
		},
	}

	for _, c := range deleteCases {
		t.Run(c.caseName, func(t *testing.T) {
			err := Delete(c.obj, c.name)
			if c.expectedErr != "" {
				if err == nil || err.Error() != c.expectedErr {
					t.Errorf("expected error %s, but was %+v", c.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Error(err)
				}
			}
			if c.result != nil {
				actual := fmt.Sprintf("%+v", c.result(c.obj))
				if actual != c.expected {
					t.Errorf("expected %s, but was %s", c.expected, actual)
				}
			}
		})
	}
}

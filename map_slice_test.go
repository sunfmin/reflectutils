package reflectutils

import (
	"testing"
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

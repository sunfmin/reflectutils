package reflectutils

import (
	"fmt"
	"testing"
)

type Person struct {
	Name        string
	Score       float64
	Gender      int
	Company     *Company
	Departments []*Department
	Projects    []*Project
	Phones      map[string]string
	Languages   map[string]Language
}

type Language struct {
	Code string
	Name string
}

type Company struct {
	Name  string
	Phone *Phone
}

type Department struct {
	Id   int
	Name string
}

type Project struct {
	Id      string
	Name    string
	Members []*Person
}

type Phone struct {
	Number string
}

type setTest struct {
	Name   string
	Value  string
	Getter func(p *Person) string
}

var cases = []*setTest{
	{
		Name:  ".Name",
		Value: "Felix",
		Getter: func(p *Person) string {
			return p.Name
		},
	},
	{
		Name:  ".Score",
		Value: "8.22",
		Getter: func(p *Person) string {
			return fmt.Sprintf("%.2f", p.Score)
		},
	},
	{
		Name:  ".Gender",
		Value: "1",
		Getter: func(p *Person) string {
			return fmt.Sprintf("%d", p.Gender)
		},
	},
	{
		Name:  ".Company.Name",
		Value: "The Plant",
		Getter: func(p *Person) string {
			return p.Company.Name
		},
	},
	{
		Name:  ".Phones.Home",
		Value: "111111111",
		Getter: func(p *Person) string {
			return p.Phones["Home"]
		},
	},
	{
		Name:  ".Phones.Company",
		Value: "2222222222",
		Getter: func(p *Person) string {
			return p.Phones["Company"]
		},
	},
	{
		Name:  ".Projects[0].Id",
		Value: "1",
		Getter: func(p *Person) string {
			return p.Projects[0].Id
		},
	},
	{
		Name:  ".Projects[1].Id",
		Value: "2",
		Getter: func(p *Person) string {
			return p.Projects[1].Id
		},
	},
	{
		Name:  ".Departments[2].Id",
		Value: "2",
		Getter: func(p *Person) string {
			return fmt.Sprintf("%d", p.Departments[2].Id)
		},
	},
	{
		Name:  ".Departments[0].Id",
		Value: "0",
		Getter: func(p *Person) string {
			return fmt.Sprintf("%d", p.Departments[0].Id)
		},
	},
	{
		Name:  ".Projects[0].Members[].Name",
		Value: "Juice",
		Getter: func(p *Person) string {
			return p.Projects[0].Members[0].Name
		},
	},
	{
		Name:  ".Projects[0].Members[2].Name",
		Value: "Bin",
		Getter: func(p *Person) string {
			return p.Projects[0].Members[2].Name
		},
	},
	{
		Name:  ".Languages.EN.Name",
		Value: "English",
		Getter: func(p *Person) string {
			return p.Languages["EN"].Name
		},
	},
}

func TestSetTheNil(t *testing.T) {
	var v *Person
	for _, c := range cases {
		err := Set(&v, c.Name, c.Value)
		if err != nil {
			t.Error(err)
			continue
		}

		val := c.Getter(v)
		if c.Value != val {
			t.Errorf("expected is %v, but was %v", c.Value, val)
		}
	}
}

func TestSetNotNil(t *testing.T) {
	v := &Person{
		Name: "F",
		Projects: []*Project{
			{Id: "3", Name: "Sendgrid"},
		},
	}

	for _, c := range cases {
		Set(&v, c.Name, c.Value)
		if c.Value != c.Getter(v) {
			t.Errorf("expected is %v, but was %v", c.Value, c.Getter(v))
		}
	}
	if v.Projects[0].Name != "Sendgrid" {
		t.Error("value was overwriten.")
	}

}

func TestSetStruct(t *testing.T) {
	var v Person

	for _, c := range cases {
		Set(&v, c.Name, c.Value)
		if c.Value != c.Getter(&v) {
			t.Errorf("expected is %v, but was %v", c.Value, c.Getter(&v))
		}
	}
}

func TestSetOtherPointers(t *testing.T) {
	var v *Person
	Set(&v, ".Company", &Company{
		Name: "The Plant",
	})

	if v.Company.Name != "The Plant" {
		t.Errorf("set failed %+v", v)
	}
}

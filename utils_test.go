package reflectutils

import (
	"fmt"
	"log"
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
	PhoneCalls  map[*Person]int
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
	// {
	// 	Name:  ".Name",
	// 	Value: "Felix",
	// 	Getter: func(p *Person) string {
	// 		return p.Name
	// 	},
	// },
	// {
	// 	Name:  ".Score",
	// 	Value: "8.22",
	// 	Getter: func(p *Person) string {
	// 		return fmt.Sprintf("%.2f", p.Score)
	// 	},
	// },
	// {
	// 	Name:  ".Gender",
	// 	Value: "1",
	// 	Getter: func(p *Person) string {
	// 		return fmt.Sprintf("%d", p.Gender)
	// 	},
	// },
	// {
	// 	Name:  ".Company.Name",
	// 	Value: "The Plant",
	// 	Getter: func(p *Person) string {
	// 		return p.Company.Name
	// 	},
	// },
	// {
	// 	Name:  ".Phones.Home",
	// 	Value: "111111111",
	// 	Getter: func(p *Person) string {
	// 		return p.Phones["Home"]
	// 	},
	// },
	// {
	// 	Name:  ".Phones.Company",
	// 	Value: "2222222222",
	// 	Getter: func(p *Person) string {
	// 		return p.Phones["Company"]
	// 	},
	// },
	// {
	{
		Name:  ".Projects[].Id",
		Value: "1",
		Getter: func(p *Person) string {
			return p.Projects[0].Id
		},
	},
	{
		Name:  ".Projects[].Id",
		Value: "2",
		Getter: func(p *Person) string {
			return p.Projects[1].Id
		},
	},
	// 	Name:  ".Departments[0].Id",
	// 	Value: "1",
	// 	Getter: func(p *Person) string {
	// 		return fmt.Sprintf("%d", p.Departments[0].Id)
	// 	},
	// },
	// {
	// 	Name:  ".Departments[1].Id",
	// 	Value: "2",
	// 	Getter: func(p *Person) string {
	// 		return fmt.Sprintf("%d", p.Departments[1].Id)
	// 	},
	// },
	// {
	// 	Name:  ".Projects[0].Members[].Name",
	// 	Value: "Juice",
	// 	Getter: func(p *Person) string {
	// 		return p.Projects[0].Members[0].Name
	// 	},
	// },
	// {
	// 	Name:  ".Projects[0].Members[2].Name",
	// 	Value: "Bin",
	// 	Getter: func(p *Person) string {
	// 		return p.Projects[0].Members[2].Name
	// 	},
	// },
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
		log.Println("=>\t", val, "\n\n\n")
		if c.Value != val {
			t.Errorf("expected is %v, but was %v", c.Value, val)
		}
	}
	fmt.Println(v)
}

func TestSetNotNil(t *testing.T) {
	v := &Person{
		Name: "F",
		Projects: []*Project{
			{Id: "3"},
		},
	}

	for _, c := range cases {
		Set(&v, c.Name, c.Value)
		if c.Value != c.Getter(v) {
			t.Errorf("expected is %v, but was %v", c.Value, c.Getter(v))
		}
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

package reflectutils

import (
	"encoding/json"
	"fmt"
)

// By given these structs
func ExampleSet_0init() {
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
}

// For How to set simple field
func ExampleSet_1setfield() {
	var p *Person
	Set(&p, "Name", "Felix")
	Set(&p, "Score", 66.88)
	Set(&p, "Gender", 1)
	printJsonV(p)
	//Output:
	// {
	// 	"Name": "Felix",
	// 	"Score": 66.88,
	// 	"Gender": 1,
	// 	"Company": null,
	// 	"Departments": null,
	// 	"Projects": null,
	// 	"Phones": null,
	// 	"Languages": null
	// }
}

// For how to set a struct property
func ExampleSet_2setstructproperty() {
	p := &Person{}
	Set(p, ".Company.Name", "The Plant")
	Set(p, ".Company.Phone.Number", "911")
	printJsonV(p)
	//Output:
	// {
	// 	"Name": "",
	// 	"Score": 0,
	// 	"Gender": 0,
	// 	"Company": {
	// 		"Name": "The Plant",
	// 		"Phone": {
	// 			"Number": "911"
	// 		}
	// 	},
	// 	"Departments": null,
	// 	"Projects": null,
	// 	"Phones": null,
	// 	"Languages": null
	// }
}

// For how to set slice and it's property
func ExampleSet_3setsliceproperty() {
	var p *Person
	Set(&p, "Departments[0].Id", 1)
	Set(&p, "Departments[0].Name", "High Tech")

	// if you jump the index for an array, It will put nil in between
	// So there will be no index out of range error.
	Set(&p, "Departments[3].Id", 1)
	Set(&p, "Departments[3].Name", "High Tech")
	printJsonV(p)
	//Output:
	// {
	// 	"Name": "",
	// 	"Score": 0,
	// 	"Gender": 0,
	// 	"Company": null,
	// 	"Departments": [
	// 		{
	// 			"Id": 1,
	// 			"Name": "High Tech"
	// 		},
	// 		null,
	// 		null,
	// 		{
	// 			"Id": 1,
	// 			"Name": "High Tech"
	// 		}
	// 	],
	// 	"Projects": null,
	// 	"Phones": null,
	// 	"Languages": null
	// }
}

// For how to set map property
func ExampleSet_4setmapproperty() {
	var p *Person
	Set(&p, "Languages.en_US.Name", "United States")
	Set(&p, "Languages.en_US.Code", "en_US")
	Set(&p, "Languages.zh_CN.Name", "China")
	Set(&p, "Languages.zh_CN.Code", "zh_CN")

	printJsonV(p)
	//Output:
	// {
	// 	"Name": "",
	// 	"Score": 0,
	// 	"Gender": 0,
	// 	"Company": null,
	// 	"Departments": null,
	// 	"Projects": null,
	// 	"Phones": null,
	// 	"Languages": {
	// 		"en_US": {
	// 			"Code": "en_US",
	// 			"Name": "United States"
	// 		},
	// 		"zh_CN": {
	// 			"Code": "zh_CN",
	// 			"Name": "China"
	// 		}
	// 	}
	// }
}

// You can do whatever deeper you like
func ExampleSet_5setdeeper() {
	var p *Person
	Set(&p, "Projects[0].Members[0].Company.Phone.Number", "911")

	printJsonV(p)
	//Output:
	// {
	// 	"Name": "",
	// 	"Score": 0,
	// 	"Gender": 0,
	// 	"Company": null,
	// 	"Departments": null,
	// 	"Projects": [
	// 		{
	// 			"Id": "",
	// 			"Name": "",
	// 			"Members": [
	// 				{
	// 					"Name": "",
	// 					"Score": 0,
	// 					"Gender": 0,
	// 					"Company": {
	// 						"Name": "",
	// 						"Phone": {
	// 							"Number": "911"
	// 						}
	// 					},
	// 					"Departments": null,
	// 					"Projects": null,
	// 					"Phones": null,
	// 					"Languages": null
	// 				}
	// 			]
	// 		}
	// 	],
	// 	"Phones": null,
	// 	"Languages": null
	// }
}

// func ExamplePopulate_1() {
// 	type Loop struct {
// 		Name string
// 		Loop *Loop
// 	}
// 	var l *Loop
// 	Populate(&l)

// 	var p *Person
// 	Populate(&p)
// 	printJsonV(p)
// 	//Output:

// }

func printJsonV(v interface{}) {
	j, _ := json.MarshalIndent(v, "", "\t")
	fmt.Printf("\n\n%s\n", j)
}

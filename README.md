

# reflectutils

reflectutils is for setting your struct value by using string name and path that follows reflectutils rules.

## The rules

- `.Name` to set a property by field name
- `.Person.Name` to set the name of the current property
- `.Person.Addresses[0].Phone` to set an element of an array property
- `.Person.Addresses[].Name` it will create a object of address and set it's property
- `.Person.MapData.Name` it can also set value to map

## How to install


```go
go get github.com/sunfmin/reflectutils
```



By given these structs
```go
	type Person struct {
	    Name        string
	    Score       float64
	    Gender      int
	    Company     *Company
	    Departments []*Department
	    Projects    []Project
	    Phones      map[string]string
	    Languages   map[string]Language
	}
	
	type Language struct {
	    Code string
	    Name string
	}
	
	type Company struct {
	    Name   string
	    Phone  *Phone
	    Phone2 **Phone
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
```

For How to set simple field
```go
	var p *Person
	Set(&p, "Name", "Felix")
	Set(&p, "Score", 66.88)
	Set(&p, "Gender", 1)
	printJsonV(p)
	fmt.Println(MustGet(&p, "Score"))
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
	// 66.88
```

For how to set a struct property
```go
	p := &Person{}
	Set(p, ".Company.Name", "The Plant")
	Set(p, ".Company.Phone.Number", "911")
	printJsonV(p)
	fmt.Println(MustGet(&p, ".Company.Phone.Number"))
	
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
	// 911
```

For how to set slice and it's property
```go
	var p *Person
	Set(&p, "Departments[0].Id", 1)
	Set(&p, "Departments[0].Name", "High Tech")
	Set(&p, "Projects[0].Name", "UIBuilder")
	// if you jump the index for an array, It will put nil in between
	// So there will be no index out of range error.
	Set(&p, "Departments[3].Id", 1)
	Set(&p, "Departments[3].Name", "High Tech")
	printJsonV(p)
	
	fmt.Println(MustGet(&p, "Departments[3].Name"))
	fmt.Println(MustGet(&p, "Departments[4].Name"))
	fmt.Println(MustGet(&p, "Projects[0].Name"))
	
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
	// 	"Projects": [
	// 		{
	// 			"Id": "",
	// 			"Name": "UIBuilder",
	// 			"Members": null
	// 		}
	// 	],
	// 	"Phones": null,
	// 	"Languages": null
	// }
	// High Tech
	// <nil>
	// UIBuilder
```

For how to set map property
```go
	var p *Person
	Set(&p, "Languages.en_US.Name", "United States")
	Set(&p, "Languages.en_US.Code", "en_US")
	Set(&p, "Languages.zh_CN.Name", "China")
	Set(&p, "Languages.zh_CN.Code", "zh_CN")
	
	printJsonV(p)
	fmt.Println(MustGet(&p, "Languages.zh_CN.Code"))
	
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
	// zh_CN
```

You can do whatever deeper you like
```go
	var p *Person
	Set(&p, "Projects[0].Members[0].Company.Phone.Number", "911")
	
	printJsonV(p)
	fmt.Println(MustGet(&p, "Projects[0].Members[0].Company.Phone.Number"))
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
	// 911
```

A new way to append to an array
```go
	var p *Person
	for i := 1; i < 6; i++ {
	    var d *Department
	    Set(&d, "Id", i)
	    Set(&d, "Name", fmt.Sprintf("Department of Energy %d", i))
	    Set(&p, "Departments[]", d)
	}
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
	// 			"Name": "Department of Energy 1"
	// 		},
	// 		{
	// 			"Id": 2,
	// 			"Name": "Department of Energy 2"
	// 		},
	// 		{
	// 			"Id": 3,
	// 			"Name": "Department of Energy 3"
	// 		},
	// 		{
	// 			"Id": 4,
	// 			"Name": "Department of Energy 4"
	// 		},
	// 		{
	// 			"Id": 5,
	// 			"Name": "Department of Energy 5"
	// 		}
	// 	],
	// 	"Projects": null,
	// 	"Phones": null,
	// 	"Languages": null
	// }
```

You could also set []byte data to string property, and vice versa.
And set string value to int, float
```go
	type Obj struct {
	    ByteProperty     []byte
	    StringProperty   string
	    IntValue         int
	    FloatValue       float64
	    IntValueForBytes int
	}
	var o *Obj
	Set(&o, "ByteProperty", "hello")
	Set(&o, "StringProperty", []byte{0x46, 0x65, 0x6c, 0x69, 0x78})
	Set(&o, "IntValue", "22")
	Set(&o, "IntValueForBytes", []byte{0x32, 0x32})
	Set(&o, "FloatValue", "22.88")
	fmt.Println(string(o.ByteProperty))
	fmt.Println(o.StringProperty)
	fmt.Println(o.IntValue)
	fmt.Println(o.FloatValue)
	fmt.Println(o.IntValueForBytes)
	//Output:
	//hello
	//Felix
	//22
	//22.88
	//22
```

If you set a property that don't exists, it gives you an error.
```go
	var p *Person
	err := Set(&p, "Whatever.Not.Exists", "911")
	
	fmt.Println(err)
	//Output:
	// no such field.
```

Get Type of a deep nested object
```go
	type Variant struct {
	    Name string
	}
	type Product struct {
	    Variants []*Variant
	    ByCode   map[string]*Variant
	}
	type Obj struct {
	    MainProduct *Product
	}
	
	var o *Obj
	
	fmt.Println(GetType(o, "MainProduct.Variants[0].Name"))
	fmt.Println(GetType(o, "MainProduct.Variants[0]"))
	fmt.Println(GetType(o, "MainProduct.Variants"))
	fmt.Println(GetType(o, "MainProduct"))
	fmt.Println(GetType(o, "MainProduct.ByCode.abc"))
	fmt.Println(GetType(o, "MainProduct.ByCode"))
	fmt.Println(GetType(o, "x123.ByCode"))
	fmt.Println(GetType(o, "MainProduct.ByCode.abc.NotExist"))
	//Output:
	//string
	//*reflectutils_test.Variant
	//[]*reflectutils_test.Variant
	//*reflectutils_test.Product
	//*reflectutils_test.Variant
	//map[string]*reflectutils_test.Variant
	//<nil>
	//<nil>
```




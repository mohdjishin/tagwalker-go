# **TagWalker-Go: Struct Tag Extractor**

**A Go library to extract custom struct tags recursively, supporting nested structures, slices, maps, and multiple tag keys.**

## **Features**

- **Extracts specified tags from struct fields, including nested fields.**
- **Supports pointers, slices, arrays, and maps.**
- **Handles multiple tag values separated by commas.**
- **Returns structured output with field paths and tag values.**

## **Installation**

```bash
go get github.com/mohdjishin/tagwalker-go/tagextractor
```


## **Usage**
### **Example**

```go
package main

import (
	"fmt"

	"github.com/mohdjishin/tagwalker-go/tagextractor"
)

func main() {
	type Address struct {
		City    string  `custom:"city_tag,location_tag" testtag:"test_City"`
		Country *string `custom:"country_tag,location_tag" testtag:"test_Country"`
	}

	type Person struct {
		Name    string   `custom:"name_tag,test_tag" testtag:"test_Name"`
		Age     int      `custom:"age_tag" testtag:"test_Age"`
		Address *Address `custom:"addr_tag" testtag:"test_Address"`
		Scores  []int    `custom:"scores_tag" testtag:"test_Scores"`
		Friends []Person `custom:"friends_tag" testtag:"test_Friends"`
		Data    map[string]struct {
			ID float64 `custom:"id_tag" testtag:"test_ID"`
		} `custom:"data_tag" testtag:"test_Data"`
	}

	p := Person{
		Address: &Address{City: "London"},
		Data: map[string]struct {
			ID float64 `custom:"id_tag" testtag:"test_ID"`
		}{
			"record1": {ID: 123.45},
		},
	}

	tagsExtractor := tagextractor.NewExtractor([]string{"custom", "testtag"})

	fmt.Println("Extracted Tags:")
	for _, tag := range tagsExtractor.Extract(p) {
		fmt.Println(tag)
	}

}
```


## **Output**

```bash
Extracted Tags:
Path: Name                           | Tag: custom          | Values: [name_tag test_tag]
Path: Name                           | Tag: testtag         | Values: [test_Name]
Path: Age                            | Tag: custom          | Values: [age_tag]
Path: Age                            | Tag: testtag         | Values: [test_Age]
Path: Address                        | Tag: custom          | Values: [addr_tag]
Path: Address                        | Tag: testtag         | Values: [test_Address]
Path: Scores                         | Tag: custom          | Values: [scores_tag]
Path: Scores                         | Tag: testtag         | Values: [test_Scores]
Path: Friends                        | Tag: custom          | Values: [friends_tag]
Path: Friends                        | Tag: testtag         | Values: [test_Friends]
Path: Data                           | Tag: custom          | Values: [data_tag]
Path: Data                           | Tag: testtag         | Values: [test_Data]
Path: Address.City                   | Tag: custom          | Values: [city_tag location_tag]
Path: Address.City                   | Tag: testtag         | Values: [test_City]
Path: Address.Country                | Tag: custom          | Values: [country_tag location_tag]
Path: Address.Country                | Tag: testtag         | Values: [test_Country]
Path: Data[record1].ID               | Tag: custom          | Values: [id_tag]
Path: Data[record1].ID               | Tag: testtag         | Values: [test_ID]
```
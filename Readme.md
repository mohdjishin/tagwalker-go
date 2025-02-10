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

	tags := tagextractor.ExtractCustomTags(p, []string{"custom", "testtag"})

	fmt.Println("Extracted Tags:")
	for _, tag := range tags {
		fmt.Printf("%-20s [%s] -> %v\n", tag.FieldPath, tag.TagKey, tag.TagValue)
	}
}
```


## **Output**

```bash
Extracted Tags:
Name                 [custom] -> [name_tag test_tag]
Name                 [testtag] -> [test_Name]
Age                  [custom] -> [age_tag]
Age                  [testtag] -> [test_Age]
Address              [custom] -> [addr_tag]
Address              [testtag] -> [test_Address]
Address.City         [custom] -> [city_tag location_tag]
Address.City         [testtag] -> [test_City]
Data                 [custom] -> [data_tag]
Data                 [testtag] -> [test_Data]
Data[record1].ID     [custom] -> [id_tag]
Data[record1].ID     [testtag] -> [test_ID]
```
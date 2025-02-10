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

	tagsExtractor := tagextractor.NewTagExtractor([]string{"custom", "testtag"})

	fmt.Println("Extracted Tags:")
	for _, tag := range tagsExtractor.Extract(p) {
		fmt.Printf("%-25s [%s] -> %v\n", tag.FieldPath, tag.TagKey, tag.TagValue)
	}
}

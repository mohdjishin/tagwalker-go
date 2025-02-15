package main

import (
	"fmt"

	"github.com/mohdjishin/tagwalker-go/tagextractor"
)

type Address struct {
	City    string `custom:"city_tag,location_tag" testtag:"test_City"`
	Country string `custom:"country_tag,location_tag" testtag:"test_Country"`
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

func main() {
	p := Person{
		Name: "John Doe",
		Age:  30,
		Address: &Address{
			City:    "New York",
			Country: "USA",
		},
		Scores: []int{1, 2, 3, 4, 5},
		Friends: []Person{
			{
				Name: "Jane Doe",
				Age:  25,
				Address: &Address{
					City:    "Los Angeles",
					Country: "USA",
				},
				Scores: []int{5, 4, 3, 2, 1},
			},
		},
		Data: map[string]struct {
			ID float64 `custom:"id_tag" testtag:"test_ID"`
		}{
			"key1": {ID: 1},
			"key2": {ID: 2},
		},
	}

	extractor := tagextractor.NewExtractor(([]string{"custom"}))
	tags := extractor.Extract(p)

	for _, tag := range tags {
		fmt.Println(tag)
	}
}

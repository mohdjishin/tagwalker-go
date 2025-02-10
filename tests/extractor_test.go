package tests

import (
	"reflect"
	"testing"

	"github.com/mohdjishin/tagwalker-go/tagextractor"
)

func TestExtractCustomTags(t *testing.T) {
	type Nested struct {
		InnerField string `custom:"inner_tag"`
	}

	type Sample struct {
		Name     string  `custom:"name_tag"`
		Age      int     `custom:"age_tag"`
		Location *string `custom:"location_tag"`
		Nested   Nested  `custom:"nested_tag"`
	}

	location := "New York"
	sample := Sample{
		Name:     "Alice",
		Age:      30,
		Location: &location,
		Nested:   Nested{InnerField: "value"},
	}

	expected := []tagextractor.FieldTag{
		{Path: "Name", Key: "custom", Values: []string{"name_tag"}},
		{Path: "Age", Key: "custom", Values: []string{"age_tag"}},
		{Path: "Location", Key: "custom", Values: []string{"location_tag"}},
		{Path: "Nested", Key: "custom", Values: []string{"nested_tag"}},
		{Path: "Nested.InnerField", Key: "custom", Values: []string{"inner_tag"}},
	}

	tagExtractror := tagextractor.NewExtractor([]string{"custom"})
	result := tagExtractror.Extract(sample)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

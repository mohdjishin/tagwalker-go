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

	expected := []tagextractor.CustomTag{
		{FieldPath: "Name", TagKey: "custom", TagValue: []string{"name_tag"}},
		{FieldPath: "Age", TagKey: "custom", TagValue: []string{"age_tag"}},
		{FieldPath: "Location", TagKey: "custom", TagValue: []string{"location_tag"}},
		{FieldPath: "Nested", TagKey: "custom", TagValue: []string{"nested_tag"}},
		{FieldPath: "Nested.InnerField", TagKey: "custom", TagValue: []string{"inner_tag"}},
	}

	result := tagextractor.ExtractCustomTags(sample, []string{"custom"})

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

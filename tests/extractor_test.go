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
		{Path: "Name", Key: "custom", Values: []string{"name_tag"}, Value: "Alice"},
		{Path: "Age", Key: "custom", Values: []string{"age_tag"}, Value: 30},
		{Path: "Location", Key: "custom", Values: []string{"location_tag"}, Value: "New York"}, // Handle pointer case
		{Path: "Nested", Key: "custom", Values: []string{"nested_tag"}, Value: sample.Nested},
		{Path: "Nested.InnerField", Key: "custom", Values: []string{"inner_tag"}, Value: "value"},
	}

	tagExtractor := tagextractor.NewExtractor([]string{"custom"})
	result := tagExtractor.Extract(sample)

	if len(result) != len(expected) {
		t.Fatalf("Expected %d fields, but got %d", len(expected), len(result))
	}

	for i := range expected {
		actualValue := result[i].Value
		expectedValue := expected[i].Value

		if v, ok := actualValue.(*string); ok {
			actualValue = *v
		}

		if result[i].Path != expected[i].Path || result[i].Key != expected[i].Key ||
			!reflect.DeepEqual(result[i].Values, expected[i].Values) ||
			!reflect.DeepEqual(actualValue, expectedValue) {
			t.Errorf("Mismatch at index %d:\nExpected: %+v\nGot: %+v", i, expected[i], result[i])
		}
	}
}

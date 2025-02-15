package tagextractor

import (
	"fmt"
	"reflect"
	"strings"
)

type FieldTag struct {
	Path   string
	Key    string
	Values []string
	Value  interface{}
}

func (f FieldTag) String() string {
	return fmt.Sprintf(
		"Path: %-35s | Tag: %-12s | Values: %-30s | Value: %v",
		f.Path, f.Key, strings.Join(f.Values, ", "), f.Value,
	)
}

func (f FieldTag) TagValue() []string {
	return f.Values
}

func (f FieldTag) TagKey() string {
	return f.Key
}

func (f FieldTag) FieldPath() string {
	return f.Path
}

type Extractor struct {
	tagKeys []string
}

func NewExtractor(tagKeys []string) *Extractor {
	return &Extractor{tagKeys: tagKeys}
}

func (e *Extractor) Extract(target interface{}) []FieldTag {
	var extractedTags []FieldTag
	rv := reflect.ValueOf(target)
	rt := reflect.TypeOf(target)

	for rt.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return extractedTags
		}
		rv = rv.Elem()
		rt = rt.Elem()
	}

	if rt.Kind() != reflect.Struct {
		return extractedTags
	}

	type fieldNode struct {
		Type  reflect.Type
		Value reflect.Value
		Path  string
	}

	queue := []fieldNode{{rt, rv, ""}}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if !node.Value.IsValid() {
			continue
		}

		for i := 0; i < node.Type.NumField(); i++ {
			field := node.Type.Field(i)
			fieldValue := node.Value.Field(i)
			fieldPath := field.Name
			if node.Path != "" {
				fieldPath = node.Path + "." + field.Name
			}

			if field.PkgPath != "" {
				continue
			}

			for _, tagKey := range e.tagKeys {
				if tagValue := field.Tag.Get(tagKey); tagValue != "" {
					extractedTags = append(extractedTags, FieldTag{
						Path:   fieldPath,
						Key:    tagKey,
						Values: strings.Split(strings.TrimSpace(tagValue), ","),
						Value:  fieldValue.Interface(),
					})
				}
			}

			ft, fv := field.Type, fieldValue
			for ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
				if fv.IsValid() && !fv.IsNil() {
					fv = fv.Elem()
				}
			}

			switch ft.Kind() {
			case reflect.Struct:
				queue = append(queue, fieldNode{ft, fv, fieldPath})

			case reflect.Slice, reflect.Array:
				if elemType := ft.Elem(); elemType.Kind() == reflect.Struct {
					for j := 0; j < fv.Len(); j++ {
						queue = append(queue, fieldNode{elemType, fv.Index(j), fmt.Sprintf("%s[%d]", fieldPath, j)})
					}
				}

			case reflect.Map:
				for _, key := range fv.MapKeys() {
					queue = append(queue, fieldNode{ft.Elem(), fv.MapIndex(key), fmt.Sprintf("%s[%v]", fieldPath, key)})
				}
			}
		}
	}
	return extractedTags
}

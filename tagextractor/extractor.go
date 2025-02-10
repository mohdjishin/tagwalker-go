package tagextractor

import (
	"fmt"
	"reflect"
	"strings"
)

type CustomTag struct {
	FieldPath string
	TagKey    string
	TagValue  []string
}

func ExtractCustomTags(target interface{}, tagKeys []string) []CustomTag {
	var tags []CustomTag
	rv := reflect.ValueOf(target)
	rt := reflect.TypeOf(target)

	for rt.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return tags
		}
		rv = rv.Elem()
		rt = rt.Elem()
	}

	if rt.Kind() != reflect.Struct {
		return tags
	}

	queue := []struct {
		Type  reflect.Type
		Value reflect.Value
		Path  string
	}{{rt, rv, ""}}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if !current.Value.IsValid() {
			continue
		}

		for i := 0; i < current.Type.NumField(); i++ {
			field := current.Type.Field(i)
			fieldValue := current.Value.Field(i)
			fieldPath := field.Name

			if current.Path != "" {
				fieldPath = current.Path + "." + field.Name
			}

			// Extract multiple tags
			for _, tagKey := range tagKeys {
				if tagValue := field.Tag.Get(tagKey); tagValue != "" {
					tags = append(tags, CustomTag{
						FieldPath: fieldPath,
						TagKey:    tagKey,
						TagValue:  strings.Split(strings.TrimSpace(tagValue), ","),
					})
				}
			}

			// Handle nested structures
			ft := field.Type
			fv := fieldValue

			for ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
				if fv.IsValid() && !fv.IsNil() {
					fv = fv.Elem()
				}
			}

			switch ft.Kind() {
			case reflect.Struct:
				queue = append(queue, struct {
					Type  reflect.Type
					Value reflect.Value
					Path  string
				}{ft, fv, fieldPath})

			case reflect.Slice, reflect.Array:
				elemType := ft.Elem()
				for elemType.Kind() == reflect.Ptr {
					elemType = elemType.Elem()
				}
				if elemType.Kind() == reflect.Struct {
					for j := 0; j < fv.Len(); j++ {
						queue = append(queue, struct {
							Type  reflect.Type
							Value reflect.Value
							Path  string
						}{elemType, fv.Index(j), fmt.Sprintf("%s[%d]", fieldPath, j)})
					}
				}

			case reflect.Map:
				for _, key := range fv.MapKeys() {
					value := fv.MapIndex(key)
					queue = append(queue, struct {
						Type  reflect.Type
						Value reflect.Value
						Path  string
					}{ft.Elem(), value, fmt.Sprintf("%s[%v]", fieldPath, key)})
				}
			}
		}
	}
	return tags
}

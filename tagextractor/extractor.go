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

type TagExtractor struct {
	tagKeys []string
}

func NewTagExtractor(tagKeys []string) *TagExtractor {
	return &TagExtractor{tagKeys: tagKeys}
}

func (te *TagExtractor) Extract(target interface{}) []CustomTag {
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

	type queueItem struct {
		Type  reflect.Type
		Value reflect.Value
		Path  string
	}

	queue := []queueItem{{rt, rv, ""}}

	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]

		if !item.Value.IsValid() {
			continue
		}

		for i := 0; i < item.Type.NumField(); i++ {
			field := item.Type.Field(i)
			fieldValue := item.Value.Field(i)
			fieldPath := field.Name
			if item.Path != "" {
				fieldPath = item.Path + "." + field.Name
			}

			for _, tagKey := range te.tagKeys {
				if tagValue := field.Tag.Get(tagKey); tagValue != "" {
					tags = append(tags, CustomTag{
						FieldPath: fieldPath,
						TagKey:    tagKey,
						TagValue:  strings.Split(strings.TrimSpace(tagValue), ","),
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
				queue = append(queue, queueItem{ft, fv, fieldPath})

			case reflect.Slice, reflect.Array:
				if elemType := ft.Elem(); elemType.Kind() == reflect.Struct {
					for j := 0; j < fv.Len(); j++ {
						queue = append(queue, queueItem{elemType, fv.Index(j), fmt.Sprintf("%s[%d]", fieldPath, j)})
					}
				}

			case reflect.Map:
				for _, key := range fv.MapKeys() {
					queue = append(queue, queueItem{ft.Elem(), fv.MapIndex(key), fmt.Sprintf("%s[%v]", fieldPath, key)})
				}
			}
		}
	}
	return tags
}

package httputil

import (
	"fmt"
	"reflect"
)

type reflectField struct {
	name       string
	value      interface{}
	field      reflect.Value
	tagOptions tagOptions
}

func reflectObject(ptr interface{}) []reflectField {
	fields := make([]reflectField, 0)

	v := reflect.ValueOf(ptr).Elem() // the struct variable
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		tag := fieldInfo.Tag           // a reflect.StructTag
		name := tag.Get("http")
		tagName, tagOptions := parseTag(name)

		if tagName == "" { // ignore if tag name unset
			continue
		}
		field := v.Field(i)

		fieldObj := reflectField{
			name:       tagName,
			value:      field.Interface(),
			field:      field,
			tagOptions: tagOptions,
		}
		fields = append(fields, fieldObj)
	}

	return fields
}

// ValidateRequireField 验证结构体字段required字段是否都赋值？通过返回nil，否则返回具体错误
func ValidateRequireField(ptr interface{}) error {

	fields := reflectObject(ptr)
	for _, obj := range fields {
		if obj.tagOptions.Contains("omit") || obj.tagOptions.Contains("omitempty") {
			continue
		}
		if obj.field.IsNil() {
			return fmt.Errorf("ValidateRequireField: %s unset", obj.name)
		}

	}
	return nil
}

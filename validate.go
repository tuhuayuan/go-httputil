package httputil

import (
	"fmt"
	"reflect"
	"strings"
)

// ValidateErrors 验证错误列表,长度为0表示没有错误
type ValidateErrors []validateError

// validateError 验证错误
type validateError struct {
	FieldNames []string `json:"fields,omitempty"`
	Code       int      `json:"code,omitempty"`
	Message    string   `json:"message,omitempty"`
}

// Add 添加一个新错误到列表
func (e *ValidateErrors) Add(fieldNames []string, code int, message string) {
	*e = append(*e, validateError{
		FieldNames: fieldNames,
		Code:       code,
		Message:    message,
	})
}

// Len 返回错误数量
func (e *ValidateErrors) Len() int {
	return len(*e)
}

// Error 实现标准错误接口
func (e *ValidateErrors) Error() string {
	return e.String()
}

// String 实现strng接口方便调试
func (e *ValidateErrors) String() string {
	var prnStr string
	for i := 0; i < len(*e); i++ {
		prnStr += fmt.Sprintf(
			"Fields: %v, code: %d, msg: %s\n",
			(*e)[i].FieldNames,
			(*e)[i].Code,
			(*e)[i].Message,
		)
	}
	return prnStr
}

// Validate 目前只支持整形非0值，字符串非空，指针非nil，slice非空
func Validate(ptrStuct interface{}) error {
	structType := reflect.TypeOf(ptrStuct)
	if structType.Kind() != reflect.Ptr || structType.Elem().Kind() != reflect.Struct {
		return nil
	}

	structType = structType.Elem()
	structValue := reflect.ValueOf(ptrStuct).Elem()

	validateErr := &ValidateErrors{}

	for i := 0; i < structType.NumField(); i++ {
		fieldType := structType.Field(i)
		fieldValue := structValue.Field(i)
		tag, ok := fieldType.Tag.Lookup("validate")
		if !ok {
			continue
		}
		if strings.Contains(tag, "required") {
			found := true
			switch fieldType.Type.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if fieldValue.Int() == 0 {
					found = false
				}
			case reflect.String:
				if fieldValue.String() == "" {
					found = false
				}
			case reflect.Ptr:
				if fieldValue.IsNil() {
					found = false
				}
			case reflect.Slice:
				if fieldValue.Len() == 0 {
					found = false
				}
			case reflect.Float32, reflect.Float64:
				if fieldValue.Float() == 0.0 {
					found = false
				}
			}

			if !found {
				validateErr.Add([]string{fieldType.Name}, -1, "required not found")
			}
		}

	}
	if validateErr.Len() == 0 {
		return nil
	}
	return validateErr
}

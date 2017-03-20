package httputil

import (
	"fmt"
)

// ValidateErrors 验证错误列表,长度为0表示没有错误
type ValidateErrors []ValidateError

// ValidateError 验证错误
type ValidateError struct {
	FieldNames []string `json:"fields,omitempty"`
	Code       int      `json:"code,omitempty"`
	Message    string   `json:"message,omitempty"`
}

// Add 添加一个新错误到列表
func (e *ValidateErrors) Add(fieldNames []string, code int, message string) {
	*e = append(*e, ValidateError{
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

package httputil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateErrors(t *testing.T) {
	var errs ValidateErrors
	errs.Add([]string{"Password"}, 1, "too simple")
	errs.Add([]string{"Name"}, 2, "too short")

	assert.Len(t, errs, 2)
	fmt.Println(&errs)
}

func TestValidate(t *testing.T) {
	// TODO 在其他项目测试过，需要完善所有支持类型的测试
}

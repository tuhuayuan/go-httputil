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

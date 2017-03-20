package middlewares

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"
)

var (
	// 模拟一个jquery.ajax数据
	postData = map[string][]string{
		"id":     []string{"tuhauyuan"},
		"title":  []string{"123456"},
		"tags[]": []string{"ruby", "jason", "mike"},
	}
)

type TodoAuthor struct {
	Name string
}

type TodoItem struct {
	ID    int64    `form:"id"`
	Title *string  `form:"title"`
	Done  bool     `form:"done"`
	Tags  []string `form:"tags,optional"`

	author *TodoAuthor
	//*string `form:""`
}

func mapForm(formStruct reflect.Value, form url.Values) {
	if formStruct.Kind() == reflect.Ptr {
		formStruct = formStruct.Elem()
	}
	formStructType := formStruct.Type()
	for i := 0; i < formStructType.NumField(); i++ {
		typeField := formStructType.Field(i)
		valueField := formStruct.Field(i)

		if typeField.Type.Kind() == reflect.Ptr && typeField.Anonymous {
			fmt.Println(typeField, valueField)
			v := reflect.New(typeField.Type.Elem())
			valueField.Set(v)

		} else if typeField.Type.Kind() == reflect.Struct {
			fmt.Println(typeField, valueField)
		} else if mapTag, ok := typeField.Tag.Lookup("form"); ok {
			fmt.Println("do ", mapTag)
		} else {
			fmt.Println("unsupport")
		}
	}

}

func TestMapForm(t *testing.T) {
	todoItem := &TodoItem{}
	mapForm(reflect.ValueOf(todoItem), postData)
	fmt.Println(*todoItem.author)
}

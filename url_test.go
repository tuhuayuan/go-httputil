package httputil

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	// 模拟一个jquery.ajax数据
	postData = map[string][]string{
		"id":       []string{"10068"},
		"name":     []string{"涂涂"},
		"title":    []string{"123456"},
		"tags[]":   []string{"ruby", "jason", "mike"},
		"ids[]":    []string{"1", "2", "3"},
		"coauthor": []string{"🐦", "🐟", "🐘"},
	}
)

type TodoAuthor struct {
	Name string `form:"name"`
}

type TodoItem struct {
	ID    int64    `form:"id"`
	Title string   `form:"title"`
	Done  bool     `form:"done,omit"`
	Tags  []string `form:"tags[],omit"`
	IDs   []*int   `form:"ids[]"`

	TodoAuthor

	//*string `form:""`
}

func TestUnpackURLValues(t *testing.T) {
	f1 := &TodoItem{}
	err := UnpackURLValues(postData, f1)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(f1)
}

func TestUnpackRequest(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		f1 := &TodoItem{}
		err := UnpackRequest(r, f1)
		if err != nil {
			w.WriteHeader(400)
			io.WriteString(w, err.Error())
		} else {
			io.WriteString(w, fmt.Sprintf("%v", f1))
		}
	}
	req := httptest.NewRequest("GET", "http://example.com/?id=1&title=hello&name=tuhuayuan&ids[]=1&ids[]=2&ids[]=3", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	assert.Equal(t, 200, w.Code)
	all, _ := ioutil.ReadAll(w.Body)
	fmt.Println(string(all))
}

func TestMapURLValues(t *testing.T) {
	f1 := &TodoItem{}
	e := &ValidateErrors{}

	mapURLValues(reflect.ValueOf(f1), postData, e)

	fmt.Println(f1)
}

func TestSetProperType(t *testing.T) {
	f1 := &TodoItem{}
	v := reflect.ValueOf(f1).Elem()

	err := setProperType(reflect.TypeOf(f1.Title), v.Field(1), "123")
	assert.NoError(t, err)
	assert.Equal(t, "123", f1.Title)
}

func TestUnpackError(t *testing.T) {
	f1 := &TodoItem{}
	e := &ValidateErrors{}
	mapURLValues(reflect.ValueOf(f1), map[string][]string{
		"id": []string{"abc"},
	}, e)
	fmt.Println(e)
}

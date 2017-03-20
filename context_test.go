package httputil

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// KeyValueMid 中间件：任意reflect.Value
func KeyValueMid(key interface{}, value interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := GetContext(r)
		ctx.Set(key, value)
		ctx.Next()
	}
}

func Test_KeyValueMid(t *testing.T) {
	type Data1 struct {
		d1 string
	}

	ctx := NewHTTPContext()
	ctx.Use(KeyValueMid("name", "tuhuayuan"))
	ctx.Use(KeyValueMid("data", &Data1{
		d1: "ruby",
	}))

	assert.HTTPBodyContains(t, ctx.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		c := r.Context()
		assert.Equal(t, c.Value("name"), "tuhuayuan")
		w.Write([]byte("ok"))
	}), "GET", "/", map[string][]string{}, "ok")
}

package httputil

import (
	"fmt"
	"net/http"
	"testing"

	"context"

	"github.com/stretchr/testify/assert"
)

func withKeyValue(key interface{}, value interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		WithValue(ctx, key, value)
		Next(ctx)
	}
}

func Test_KeyValueMid(t *testing.T) {
	ctx := WithHTTPContext(context.Background())
	Use(ctx, withKeyValue("name", "tuhuayuan"))

	assert.HTTPBodyContains(t, HandleFunc(ctx,
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			ctx := r.Context()
			fmt.Println(ctx.Value("name"))
			w.Write([]byte("ok"))

		}), "GET", "/", map[string][]string{}, "ok")
}

func TestContextNil(t *testing.T) {
	ctx := WithHTTPContext(nil)
	Use(ctx, withKeyValue("name", "tuhuayuan"))

	assert.HTTPBodyContains(t, HandleFunc(ctx,
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			ctx := r.Context()
			fmt.Println(ctx.Value("name"))
			w.Write([]byte("ok"))

		}), "GET", "/", map[string][]string{}, "ok")
}

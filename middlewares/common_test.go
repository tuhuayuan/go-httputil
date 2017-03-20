package middlewares

import (
	"net/http"
	"testing"
	"zonst/qipai-golang-libs/httputil"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_Logger(t *testing.T) {
	ctx := httputil.NewHTTPContext()
	ctx.Use(Logger(logrus.New()))

	assert.HTTPBodyContains(t, ctx.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}), "GET", "/", map[string][]string{}, "ok")
}

func Test_KeyValueMid(t *testing.T) {
	type Data1 struct {
		d1 string
	}

	ctx := httputil.NewHTTPContext()
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

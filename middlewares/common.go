package middlewares

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"

	"zonst/qipai-golang-libs/httputil"
)

// Logger 请求日志中间件
func Logger(logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := httputil.GetContext(r)
		start := time.Now()
		ctx.Next()
		logger.Infof("[%s] [%s] from [%s] used %.3fsecs",
			r.Method, r.URL.Path, r.RemoteAddr, time.Now().Sub(start).Seconds())
	}
}

// KeyValue 存储任意KV中间件
func KeyValue(key interface{}, value interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := httputil.GetContext(r)
		ctx.Set(key, value)
		ctx.Next()
	}
}

// BindForm 绑定请求参数中间件
func BindForm(ptrStruct interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

// ErrorHandler 错误处理中间件
func ErrorHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

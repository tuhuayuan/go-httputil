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

// KeyValueMid 中间件：任意KV
func KeyValueMid(key interface{}, value interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := httputil.GetContext(r)
		ctx.Set(key, value)
		ctx.Next()
	}
}

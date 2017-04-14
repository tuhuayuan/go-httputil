package httputil

import (
	"container/list"
	"context"
	"net/http"
	"reflect"
)

type httpContextKey string

const (
	httpContextName httpContextKey = "_httpcontext_"
)

// Context 上下文
type httpContext struct {
	ele  *list.Element
	head *list.List
	w    http.ResponseWriter
	r    *http.Request
}

// WithHTTPContext 新建一个Context
func WithHTTPContext(parent context.Context) context.Context {
	return context.WithValue(parent, httpContextName, &httpContext{
		head: list.New(),
	})
}

// Use 设置一个中间件
func Use(ctx context.Context, f http.HandlerFunc) context.Context {
	httpCtx := ctx.Value(httpContextName).(*httpContext)
	httpCtx.head.PushBack(f)
	return ctx
}

// HandleFunc 包装http.HandlerFunc
func HandleFunc(ctx context.Context, handlers ...http.HandlerFunc) http.HandlerFunc {
	httpCtx := ctx.Value(httpContextName).(*httpContext)

	// 生成静态上下文
	staticHead := list.New()
	staticHead.PushBackList(httpCtx.head)
	for _, f := range handlers {
		staticHead.PushBack(f)
	}
	// 已一个nil表示结尾
	staticHead.PushBack(nil)

	return func(w http.ResponseWriter, r *http.Request) {
		rawContext := r.Context()
		// 生成动态上下文
		dynamic := &httpContext{
			head: staticHead,
			ele:  staticHead.Front(),
			w:    w,
		}
		if rawContext != nil {
			dynamic.r = r.WithContext(context.WithValue(rawContext, httpContextName, dynamic))
		} else {
			dynamic.r = r.WithContext(context.WithValue(context.Background(), httpContextName, dynamic))
		}
		Next(dynamic.r.Context())
	}
}

// Next 调用下一个中间件
func Next(ctx context.Context) {
	httpCtx := ctx.Value(httpContextName).(*httpContext)
	v := httpCtx.ele.Value
	if !reflect.ValueOf(v).IsNil() {
		handler := (httpCtx.ele.Value).(http.HandlerFunc)
		httpCtx.ele = httpCtx.ele.Next()
		handler(httpCtx.w, httpCtx.r)
	}
}

// WithValue 添加数据
func WithValue(ctx context.Context, key, val interface{}) context.Context {
	httpCtx := ctx.Value(httpContextName).(*httpContext)
	httpCtx.r = httpCtx.r.WithContext(context.WithValue(ctx, key, val))
	return httpCtx.r.Context()
}

package httputil

import (
	"container/list"
	"context"
	"net/http"
	"reflect"
	"time"
)

// HTTPContext HTTP加强版上下文实现
type HTTPContext interface {
	context.Context

	Set(key interface{}, val interface{})
	Use(f http.HandlerFunc) HTTPContext
	HandleFunc(f ...http.HandlerFunc) http.HandlerFunc
	Next()
}

// Context 上下文
type chainedContext struct {
	ctx  context.Context
	ele  *list.Element
	head *list.List
	w    http.ResponseWriter
	r    *http.Request
}

// NewHTTPContext 新建一个Context
func NewHTTPContext() HTTPContext {
	return &chainedContext{
		ctx:  context.Background(),
		head: list.New(),
	}
}

// GetContext 获取HTTPContext或者panic.
func GetContext(r *http.Request) HTTPContext {
	return r.Context().(HTTPContext)
}

func (c *chainedContext) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *chainedContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *chainedContext) Err() error {
	return c.ctx.Err()
}

func (c *chainedContext) Value(key interface{}) interface{} {
	return c.ctx.Value(key)
}

// Set 设置context自定义值
func (c *chainedContext) Set(key interface{}, value interface{}) {
	c.ctx = context.WithValue(c.ctx, key, value)
}

// Use 设置一个中间件
func (c *chainedContext) Use(f http.HandlerFunc) HTTPContext {
	c.head.PushBack(f)
	return c
}

// HandleFunc 包装http.HandlerFunc
func (c *chainedContext) HandleFunc(handlers ...http.HandlerFunc) http.HandlerFunc {
	// 生成静态上下文
	staticHead := list.New()
	staticHead.PushBackList(c.head)
	for _, f := range handlers {
		staticHead.PushBack(f)
	}
	staticHead.PushBack(nil)

	return func(w http.ResponseWriter, r *http.Request) {
		parentCtx := r.Context()
		if parentCtx == nil {
			parentCtx = context.Background()
		}
		// 生成动态上下文
		chainedCxt := &chainedContext{
			ctx:  parentCtx,
			head: staticHead,
			ele:  staticHead.Front(),
			w:    w,
		}
		chainedCxt.r = r.WithContext(chainedCxt)
		chainedCxt.Next()
	}
}

// Next 调用下一个中间件
func (c *chainedContext) Next() {
	v := c.ele.Value
	if !reflect.ValueOf(v).IsNil() {
		cur := (c.ele.Value).(http.HandlerFunc)
		c.ele = c.ele.Next()
		cur(c.w, c.r)
	}
}

package go_aop

import (
	"context"
	"math"
)

const ABORT = math.MaxInt32 - 10000

type Context struct {
	offset   int
	handlers []func(*Context)
	Ctx      context.Context
	param    []CustomizeParam
}

type CustomizeParam interface {
	// CustomizeContext 获取参数
	CustomizeContext() map[string]interface{}
}

func newContext(ctx context.Context, param ...CustomizeParam) *Context {
	return &Context{
		offset:   -1,
		handlers: make([]func(*Context), 0, 10),
		Ctx:      ctx,
		param:    param,
	}
}

func (ctx *Context) GetParam() []CustomizeParam {
	return ctx.param
}

func (ctx *Context) Next() {
	ctx.offset++
	s := len(ctx.handlers)
	for ; ctx.offset < s; ctx.offset++ {
		if !ctx.isAbort() {
			ctx.handlers[ctx.offset](ctx)
		} else {
			return
		}
	}
}

func (ctx *Context) Reset() {
	ctx.offset = -1
	ctx.handlers = ctx.handlers[:0]
}

// stop middleware chain
func (ctx *Context) Abort() {
	ctx.offset = math.MaxInt32 - 10000
}

func (ctx *Context) isAbort() bool {
	if ctx.offset >= ABORT {
		return true
	}
	return false
}

func (ctx *Context) addHandler(f func(ctx *Context)) {
	ctx.handlers = append(ctx.handlers, f)
}

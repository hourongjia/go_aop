package go_aop

import "context"

type WrapF struct {
	f   func(ctx context.Context, param ...CustomizeParam)
	ctx *Context
}

func WrapFunc(f func(ctx context.Context, param ...CustomizeParam), ctx context.Context, param ...CustomizeParam) *WrapF {
	return &WrapF{
		f:   f,
		ctx: newContext(ctx, param...),
	}
}

func (wf *WrapF) Use(f func(c *Context)) {
	wf.ctx.addHandler(f)
}

func (wf *WrapF) Handle() {
	wf.ctx.handlers = append(wf.ctx.handlers, func(c *Context) {
		wf.f(wf.ctx.Ctx, wf.ctx.param...)
	})

	if len(wf.ctx.handlers) > 0 {
		wf.ctx.Next()
	}
	wf.ctx.Reset()
}

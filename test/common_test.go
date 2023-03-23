package test

import (
	"context"
	"fmt"
	"go_aop"
	"testing"
)

var f2 = func(ctx context.Context, param ...go_aop.CustomizeParam) {
	fmt.Println("handle f2")
}

type MyCustomizeParam struct {
	Param string // 这里定义你自己需要的上下文参数，比如string或者你自己定义的action
}

func (receiver MyCustomizeParam) CustomizeContext() map[string]interface{} {
	p := map[string]interface{}{
		"hourongjia": "hourongjia",
	}
	return p
}

func limitRate(c *go_aop.Context) {
	fmt.Println("limit rate in")
	list := c.GetParam()
	list[0].CustomizeContext()
	fmt.Println(list[0])
	c.Next()
	fmt.Println("limit rate out ")
}

func AbortAdapt(c *go_aop.Context) {
	fmt.Println("AbortAdapt in")
	c.Abort()
	fmt.Println("AbortAdapt out ")
}

func AbortF(c *go_aop.Context) {
	fmt.Println("abort f in")
	c.Next()
	fmt.Println("abort f out")
}

// TestUsualUsage 测试平常用法
func TestUsualUsage(t *testing.T) {
	a := "kebi"
	var amy = MyCustomizeParam{
		Param: a,
	}

	ctx := context.Background()
	wrapF := go_aop.WrapFunc(f2, ctx, amy)
	wrapF.Use(limitRate)
	wrapF.Use(AbortF)
	wrapF.Handle()
}

func TestAbort(t *testing.T) {
	a := "kebi"
	var amy = MyCustomizeParam{
		Param: a,
	}

	ctx := context.Background()
	wrapF := go_aop.WrapFunc(f2, ctx, amy)
	wrapF.Use(limitRate)
	wrapF.Use(AbortAdapt)

	wrapF.Handle()
}

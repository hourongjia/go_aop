# 简介
go_aop 是一个go语言实现aop功能的包，主要参考gin框架中的拦截器实现，具有对核心方法拦截，过滤的功能

# 使用方法

步骤1：用户确定要`拦截`或者`过滤`的方法或者函数

``` 
var f2 = func(ctx context.Context, param ...go_aop.CustomizeParam) {
	fmt.Println("handle f2")
}
``` 

go_aop.CustomizeParam是一个接口，如果用户被拦截的方法需要传入自定义的上下文，需要实现go_aop.CustomizeParam，比如，针对maia，maia的上下文变量是maia.Action:

``` 
type MaiaCustomizeParam struct {
	action maia.Action // 这里定义你自己需要的上下文参数，比如string或者你自己定义的action
}

func (receiver MyCustomizeParam) CustomizeContext() map[string]interface{} {
	p := map[string]interface{}{
		"maiaAction": receiver.action,
	}
	return p
}

``` 

随后，用户可以在执行方法f2的方法中，获取到action对象

``` 
var f2 = func(ctx context.Context, param ...go_aop.CustomizeParam) {
    map := param[0].CustomizeContext()       
    action := map["maiaAction"]                //获取用户被拦截方法自己的上下文
    action.Do()                                //执行上下文方法
	fmt.Println("handle f2")
}

``` 


步骤2：实例化go_aop.WrapF对象

``` 
func WrapFunc(f func(ctx context.Context, param ...CustomizeParam), ctx context.Context, param ...CustomizeParam) *WrapF 
``` 

步骤3：实现`拦截`或者`过滤`函数

当用户希望实现`过滤`操作的时候，

``` 
func XXX(c *go_aop.Context) {
   fmt.Println("前置操作")
   c.Next()  //过滤功能用c.Next()方法，表示继续执行函数
   fmt.Println("后置操作")
}
``` 

当用户希望实现`拦截`操作的时候，

``` 
func XXX(c *go_aop.Context) {
   fmt.Println("前置操作")
   c.Abort()  //过滤功能用c.Abort()方法，表示不会执行下个拦截器的方法
   fmt.Println("后置操作")
}
``` 

当用户希望将`拦截器`或者`过滤器`中写入变量到上下文中，并且全局可读，可以如下demo所示：

``` 

func AFilter(c *go_aop.Context) {
   fmt.Println("A前置操作")
   c.Ctx = context.WithValue(c.Ctx,"key","value")  //写入上下文变量
   c.Next()  //过滤功能用c.Abort()方法，表示不会执行下个拦截器的方法
   fmt.Println(c.Ctx.Value("key"))                 // 同一个过滤器中读取之前写入上下文变量
   fmt.Println("A后置操作")
}

func BFilter(c *go_aop.Context) {
   fmt.Println("B前置操作")
   c.Next()  //过滤功能用c.Abort()方法，表示不会执行下个拦截器的方法
   fmt.Println(c.Ctx.Value("key"))                 // 不同过滤器读取之前写入上下文变量
   fmt.Println("B后置操作")
}


``` 

如果用户希望过滤器中的上下文变量不共享，即A过滤器只能读到A的上下文变量，B的过滤器只能读到b的上线文变量，如下操作：

```

func AFilter(c *go_aop.Context) {
   fmt.Println("A前置操作")
   ctx := c.Ctx
   ctx = context.WithValue(ctx,"key","value")  //写入上下文变量
   c.Next()  //过滤功能用c.Abort()方法，表示不会执行下个拦截器的方法
   fmt.Println(ctx.Value("key"))                 // 同一个过滤器中读取之前写入上下文变量,可以读取到
   fmt.Println("A后置操作")
}

func BFilter(c *go_aop.Context) {
   fmt.Println("B前置操作")
   c.Next()  //过滤功能用c.Abort()方法，表示不会执行下个拦截器的方法
   fmt.Println(c.Ctx.Value("key"))                 // 不同过滤器读取之前写入上下文变量，无法读取到数据
   fmt.Println("B后置操作")
} 

``` 

即不要将withValue完的context写回c.Ctx



步骤4：将步骤2的函数加入到wrap对象中

eg:
``` 
    ctx := context.Background()
	wrapF := go_aop.WrapFunc(f2, ctx)
	// 绑定limit
	wrapF.Use(limitRate)
	wrapF.Use(AbortAdapt)
``` 


步骤5：调用wrap对象的Handle方法，开始执行

eg:
``` 
    ctx := context.Background()
	wrapF := go_aop.WrapFunc(f2, ctx)
	wrapF.Use(limitRate)
	wrapF.Use(AbortAdapt)
	// 执行handler方法
	wrapF.Handle()
``` 


完成，用户可以参考源代码中的common_test.go文件，里边有详细的测试demo和使用demo，欢迎使用

package code

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func TestContext1() {
	//两个默认上下文 这两兄弟都返回一个emptyCtx

	//获取一个根节点。相当于拿到了Context的Root节点。注意在代码中不要随便使用。
	cbg := context.Background()
	//获取一个空节点。理论上可以把他当做根节点，请不要这样搞。这个方法只会用在需要传一个CTX，但实际不产生任何用处时，充当占位符。
	context.TODO()

	//常用的4种上下文

	//常用的一种方法，返回一个可以取消的上下文。当父辈执行取消操作时，会将信号同步给他的子孙们，他们同时也会响应取消当前的操作。
	_, cf := context.WithCancel(cbg)
	cf()

	//最常使用的一个方法。带有一个超时机制的上下文。可以手动触发取消也可以在超时之后自动触发取消。
	_, cf1 := context.WithTimeout(cbg, 10*time.Second)
	cf1()

	//不常用的一个方法。 带有一个时间点的上下文。可以手动触发取消也可以在到了时间点之后自动触发取消。
	h1, _ := time.ParseDuration("1h")
	_, cf2 := context.WithDeadline(cbg, time.Now().Add(h1))
	cf2()

	//最不常用的方法，通常情况下，我们是不会在上下文里穿一些数据的。最常见的只有传Trace和Token
	context.WithValue(cbg, "name", "关注香香编程喵喵喵")
	//WithValue 不会响应取消信号，但它的子孙辈如果是可以取消的，则会继续响应。
}

func TestContest2() {
	c0 := context.Background()        //ROOT
	c1, cf1 := context.WithCancel(c0) //c1辈，一个可以取消的上下文

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done(): //响应Context的取消信号
				fmt.Printf("关注香香编程喵喵喵！err:%s\n", ctx.Err())
				return
			default:
				fmt.Printf("喵师傅正在休息！\n")
			}
		}
	}(c1)
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("关注香香编程谢谢喵喵喵！err:%s\n", ctx.Err())
				return
			default:
				fmt.Printf("喵师傅正在喝茶！\n")
			}
		}
	}(c1)

	time.Sleep(1 * time.Millisecond)
	cf1() //手动调用终止函数
	time.Sleep(3 * time.Second)
	fmt.Printf("程序终止！")

}

type Context interface {
	//Deadline 返回具体的死线时间。ok代表是否设置过这个时间。
	Deadline() (deadline time.Time, ok bool)

	//Done 又是一个经典用法。返回一个单向通道，类型是空结构体。它返回的就是取消信号。
	Done() <-chan struct{}

	//Err 当Context取消时，获取具体的错误信息。超时还是手动取消。
	Err() error
	//Value 设置变量 不常用
	Value(key interface{}) interface{}
}
type canceler interface {
	cancel(removeFromParent bool, err error)
	Done() <-chan struct{}
}

// WithCancel 使用了这个结构体
type cancelCtx struct {
	Context
	mu       sync.Mutex            // 互斥锁，保证上下文并发安全
	done     atomic.Value          // 记录取消的状态
	children map[canceler]struct{} // 记录自己的子辈
	err      error                 // 记录取消时的错误信息
}

// WithTimeout 和 WithDeadline 使用了这个结构体
type timerCtx struct {
	cancelCtx
	timer    *time.Timer
	deadline time.Time
}

// WithValue 使用了这个
type valueCtx struct {
	Context
	key, val any
}

//type Context struct {
//HTTP相关字段
//writermem responseWriter //返回内容
//Request   *http.Request  //Http官方包的Request
//Writer    ResponseWriter //自定义的Response

//URL相关
//Params   Params        //请求参数，底层是一个切片
//handlers HandlersChain //方法链，一个HandlerFunc数组
//index    int8          //HandlerFunc数组的索引，也是Abort的标志
//fullPath string        //当次请求的全路径

//原始数据的指针，注意这三个值都是指针。
//engine       *Engine        //gin框架的主引擎
//params       *Params        //请求初始值，没有找到实际用途
//skippedNodes *[]skippedNode //路由中配置的需要跳过的路径节点

// mu 为Keys增加的读写锁。Map并发不安全，Context又需要并发。
//mu sync.RWMutex

// Keys 让Context能够提供存储一些数据的能力。Set()和Get()的数据就放在这里。
//Keys map[string]any

// Errors 是一个error的指针切片，用来存储使用该Context的方法和中间件所产生的Error
//Errors errorMsgs

// Accepted 用来记录该路径允许请求的类型，就是header里的Accept
//Accepted []string

// 这两兄弟就是url官方包，gin框架中的多个query的方法，都是在此基础上的封装。
//queryCache url.Values
//formCache  url.Values

//sameSite Cookie相关的配置字段，只在设置Cookie时使用。
//sameSite http.SameSite
//}

//func (c *Context) Deadline() (deadline time.Time, ok bool) {
//	// ContextWithFallback默认是false 一般都不会回退的。
//	if !c.engine.ContextWithFallback || c.Request == nil || c.Request.Context() == nil {
//		return
//	}
//	//如果允许回退，那么Context会直接回退为Request下的Context。
//	return c.Request.Context().Deadline()
//}

//省略掉 Done 和 Err 内容同上

// Value 会优先返回元数据中的Key，如果没有才会尝试去Request中的Context中查询。
//func (c *Context) Value(key any) any {
//	//两种特殊情况
//	if key == 0 {
//		return c.Request //直接返回Request结构体
//	}
//	if key == ContextKey {
//		return c //直接返回Context本身
//	}
//	if keyAsString, ok := key.(string); ok {
//		if val, exists := c.Get(keyAsString); exists {
//			return val
//		}
//	}
//	if !c.engine.ContextWithFallback || c.Request == nil || c.Request.Context() == nil {
//		return nil
//	}
//	return c.Request.Context().Value(key)
//}

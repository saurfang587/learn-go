package code

import (
	"context"
	"fmt"
	"time"
)

func TestTick0() {
	//常见的定时器方法
	<-time.Tick(time.Second)                       //不常用，参数不能小于等于0，底层调用的NewTicker，返回一个单向通道
	<-time.After(time.Second)                      //底层调用NewTimer，返回一个单向通道
	<-time.NewTicker(time.Second).C                //最常用的方法，会展开讲
	<-time.NewTimer(time.Second).C                 //在业务场景下用的少，在中间件里用到很多。
	time.AfterFunc(time.Second, func() { /*do*/ }) //底层是对timer的封装，用的也不多，个别场景下要比NewTimer好用的多。
	time.Sleep(time.Second)                        //最常用的方法。
}

func TestTick1() {
	//NewTimer只会执行一次
	t := time.NewTimer(3 * time.Second)

	for {
		select {
		case <-t.C: //timer只会触发一次chan
			fmt.Printf("关注香香编程喵喵喵，关注香香编程谢谢喵喵喵！")
		}
	}
}

func TestTick2() {
	//最常用的定时任务方案
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop() //随手关闭。如果不关闭有可能会引发内存溢出 注意，我们关闭ticker,但不会关闭底层通道。
	for {
		select {
		case t := <-ticker.C: //ticker会一直重复发送。
			fmt.Println("关注香香编程喵喵喵，关注香香编程谢谢喵喵喵！", t.Unix())
		}
	}
}

func TestTick3() {
	//错误的用法：这样子写会导致重复定义多个定时器！并且，这样子时间长的永远不会被执行。
	for {
		select {
		case t := <-time.Tick(4 * time.Second): //Tick相当于快速启动一个ticker。
			fmt.Println("关注香香编程喵喵喵，关注香香编程谢谢喵喵喵！Tick:", t.Unix())
		case a := <-time.After(2 * time.Second): //After相当于快速启动一个timer。
			fmt.Println("关注香香编程喵喵喵，关注香香编程谢谢喵喵喵！After:", a.Unix())
		}
	}
	//这个写法必然会导致不停的创建一个timer 引发内存泄露。通常用不到这两个方法。

}

func TestTick4() {
	//AfterFunc 相当于延时执行某个方法。并且只会执行一次。
	//它不需要开发人员再去管理Timer，只需要完成需要延时执行的方法即可。
	time.AfterFunc(3*time.Second, func() {
		fmt.Println("关注香香编程喵喵喵，关注香香编程谢谢喵喵喵！AfterFunc:", time.Now().Unix())
	})
	select {}
}

type f func(c context.Context) error

type MTicker struct {
	C      context.Context
	T      *time.Ticker
	Worker f
}

func NewMTicker(c context.Context, d int, f f) *MTicker {
	return &MTicker{
		C:      c,
		T:      time.NewTicker(time.Duration(d) * time.Second),
		Worker: f,
	}
}

func (t *MTicker) Start() {
	for {
		select {
		case <-t.T.C:
			if err := t.Worker(t.C); err != nil {
				//处理异常
			}
		}
	}
}

func (t *MTicker) Stop() {
	t.T.Stop()
}

func TestMain() {
	t := NewMTicker(context.Background(), 1, func(c context.Context) error {
		fmt.Printf("关注香香编程喵喵喵，关注香香编程谢谢喵喵喵！")
		return nil
	})
	defer t.Stop()
	t.Start()
}

// 一个场景
func TestTicker() {
	tick := time.NewTicker(time.Second / 100) //一秒执行一百次

	go func() {
		time.Sleep(10 * time.Second) //10秒之后结束定时器，使for range 终止
		tick.Stop()                  //注意啊 这个stop不会关闭底层的通道，避免出现错误
		fmt.Printf("tick is closed\n")
	}()
	for range tick.C { //一个不常用的用法，tick.C不会close，代码会一直阻塞在这里
		fmt.Printf("关注香香编程喵喵喵，关注香香编程谢谢喵喵喵！") //理论上会执行10*100次
	}
	fmt.Printf("TestPprof end \n") //这一行是无法被打印出来的
	return
}

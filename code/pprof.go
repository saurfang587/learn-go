package code

import (
	"fmt"
	"net/http"
	"os"
	"runtime/pprof"
	"time"

	_ "net/http/pprof"
)

func creatFile() {
	file1, _ := os.Create("./cpu.pprof") //开启一个文件
	_ = pprof.StartCPUProfile(file1)     //记录CPU的使用情况
	file2, _ := os.Create("./mem.pprof")
	_ = pprof.WriteHeapProfile(file2) //记录内存的使用情况
	defer func() {
		pprof.StopCPUProfile() //随手关闭文件
		_ = file1.Close()
		_ = file2.Close()
	}()
}

func TestPprof() {
	creatFile()
	tick := time.NewTicker(time.Second / 100) //一秒执行一百次
	closed := make(chan bool)
	go func() {
		time.Sleep(10 * time.Second) //10秒之后结束定时器，使for range 终止
		tick.Stop()
		closed <- true
		fmt.Printf("tick is closed\n")
	}()
	for range tick.C { //一个不常用的用法，tick.C只要不close，就会一直阻塞在这里
		go forSelect(closed) //理论上会执行10*100次
		if <-closed {        //收到关闭消息后，跳出循环
			break
		}
	}
	fmt.Printf("关注香香编程喵喵喵，关注香香编程谢谢喵喵喵！\n")
}

func forSelect(c chan bool) {
	var ch chan int //一个典型的错误用法，此时ch是nil
	for {
		select {
		case _ = <-ch: //读取一个为nil的通道，一定会阻塞
			fmt.Printf("关注香香编程喵喵喵，关注香香编程谢谢喵喵喵！")
		case fl := <-c: //当定时器终止时跳出循环
			if fl {
				fmt.Printf("定时器已终止\n")
				return
			}
		default:
			//每次循环都会默认走这里,可以去掉注释看下
			//fmt.Printf("default")
		}
	}
}

func TestWebPprof() {
	go func() {
		if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
			fmt.Printf("%v", err)
		}
	}()

	tick := time.NewTicker(time.Second / 100) //一秒执行一百次
	closed := make(chan bool)
	go func() {
		time.Sleep(10 * time.Second) //10秒之后结束定时器，使for range 终止
		tick.Stop()
		closed <- true
		fmt.Printf("tick is closed\n")
	}()
	for range tick.C { //一个不常用的用法，tick.C只要不close，就会一直阻塞在这里
		go forSelect(closed) //理论上会执行10*100次
		if <-closed {        //收到关闭消息后，跳出循环
			break
		}
	}
	fmt.Printf("关注香香编程喵喵喵，关注香香编程谢谢喵喵喵！\n")
}

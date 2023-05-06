package code

import (
	"fmt"
	"sync"
)

// 单例模式

var GlobalMem *Mem //一个全局变量

type Mem struct {
	name string
	Age  string //不要这样定义字段，调用方能直接修改这个字段，不安全。
}

func (m *Mem) GetName() string {
	return m.name
}

func (m *Mem) SetName() {
	//最好不要提供这种方法，不要把单例的控制权交给任何调用方。
}

func NewGlobalMem() *Mem {
	//单例的核心方法
	if GlobalMem != nil {
		return GlobalMem
	}

	GlobalMem = &Mem{name: "关注香香编程喵喵喵，关注香香编程谢谢喵喵喵！"}
	return GlobalMem
}

func TestMem() {
	//两种用法
	//1.直接在main里调用NewGlobalMem,此时全局GlobalMem已经被初始化了。
	_ = NewGlobalMem()
	name := GlobalMem.GetName()
	name = NewGlobalMem().GetName() //这两种方法都可以。
	fmt.Printf(name)

	//2.如果没有在main中初始化，后续直接调用NewGlobalMem即可。
	name = NewGlobalMem().GetName() //但是，这样子做有一个并发安全问题。如果同时10个协程使用，会创建10次。
	fmt.Printf(name)

}

var once = &sync.Once{}

func NewGlobalMemOnce() *Mem {
	once.Do(func() {
		//sync.Once 可以保证无论多少并发和调用，这个实例化只会被执行一次。
		GlobalMem = &Mem{name: "关注香香编程喵喵喵，关注香香编程谢谢喵喵喵！"}
	})
	return GlobalMem
}

var GlobalMem1 = &Mem{name: "关注香香编程喵喵喵，关注香香编程谢谢喵喵喵！"} //在声明GlobalMem1的时候，直接把它实现掉。

func TestMem1() {
	//后续直接使用即可，不需要实例化，也不存在并发安全问题。
	name := GlobalMem1.GetName()
	fmt.Printf(name)
}

//工厂模式

type Pet struct {
	Name string
	Age  string //不要这样定义字段，调用方能直接修改这个字段，不安全。
}

func (p *Pet) String() {
	fmt.Printf(p.Name)
}

func NewPet(name, age string) *Pet {
	//典型的工厂模式用法，一般没人用
	return &Pet{
		Name: name,
		Age:  age,
	}
}

func TestPet() {
	p := NewPet("关注香香编程喵喵喵，关注香香编程谢谢喵喵喵！", "1")

	//我直接声明它不香么
	p = &Pet{
		Name: "关注香香编程喵喵喵，关注香香编程谢谢喵喵喵！",
		Age:  "1",
	}
	//有些人会说你这样声明，没有工厂模式优雅，可读性不好。只能说，你不懂Go的精髓。
	//不要干那些画蛇添足的事儿。工厂模式一般不用在这里。

	p.String()
}

//简单工厂

type Config interface { //定义一个接口
	String()
}

type ConfigJson struct { //第一个接口实现
}

func (c *ConfigJson) String() {
}

func newConfigJson() *ConfigJson { //工厂模式
	return &ConfigJson{}
}

type ConfigYaml struct { //第二个接口实现
}

func (c *ConfigYaml) String() {
}

func newConfigYaml() *ConfigYaml {
	return &ConfigYaml{}
}

func TestConfig(t string) Config { //这里一定是返回一个interface
	//根据传入的参数，选择对应的实例化对象，这是简单工厂的典型用法。
	switch t {
	case "json":
		return newConfigJson()
	case "yaml":
		return newConfigYaml()
	default:
		return nil
	}
}

//观察者模式

type Observer interface { // 观察者接口，也叫订阅者
	OnCall(int) //触发响应的方法
}

type Notifier interface { // 被观察者，也叫发布者
	Register(*Observer)   //把一个观察者注册进来
	Deregister(*Observer) //移除一个观察者，有些实现中没有移除的方法，这样会不完整，缺了一点灵活性
	Notify(int)           //发出消息，通知观察者们状态变化,int 来代表一个状态
}

// NotifierOne 实现一个发布者
type NotifierOne struct {
	//observers []Observer //如果不需要移除操作，那么使用切片就好
	observers map[Observer]struct{} //需要移除的话，显然使用MAP更合适。
	//注意这里有两个细节，Key使用了interface,但是在接收的时候使用的是指针，Value是空结构体。我们之前都分享过这些细节。
	status int //被观察的状态
}

func (n *NotifierOne) Register(o Observer) {
	n.observers[o] = struct{}{} //将观察者装进Map
}

func (n *NotifierOne) Deregister(o Observer) {
	delete(n.observers, o)
}

func (n *NotifierOne) Notify(status int) {
	n.status = status
	for observer, _ := range n.observers {
		//逐个调用MAP中观察者的OnCall 方法，以此来实现通知操作。
		observer.OnCall(status)
		//这样也可以，通过并发的方式进行调用。但，任何问题只要开启并发，就会引入其他问题，需要提前想清楚
		go observer.OnCall(status)
	}
}

// Observer1 观察者1号
type Observer1 struct {
}

func (o *Observer1) OnCall(status int) {
	//不同的观察者，可以做不同的事情，也可以做相同的事情。
	fmt.Printf("关注香香编程喵喵喵！Status:%d\n", status)
}

// Observer2 观察者2号
type Observer2 struct {
}

func (o *Observer2) OnCall(status int) {
	fmt.Printf("关注香香编程谢谢喵喵喵！Status:%d\n", status)
}

func TestObs() {
	//实例化一个被观察者
	notifier := &NotifierOne{
		observers: make(map[Observer]struct{}),
		status:    0,
	}

	//实例化两个观察者
	o1 := Observer1{}
	o2 := Observer2{}
	notifier.Register(&o1)
	notifier.Register(&o2)

	fmt.Printf("observers len:%d\n", len(notifier.observers)) // 2
	//触发通知
	notifier.Notify(1)
}

// 观察者改造
type eveFun func(a *Event, cond *sync.Cond) //订阅方法

type Event struct { //利用这个结构体，作为状态的载体
	Status int
	FG     int //计数器
	sc     *sync.Cond
	fs     []eveFun //订阅方法的数组
}

func (e *Event) Run() {
	for _, f2 := range e.fs {
		go f2(e, e.sc) //让订阅方法跑起来
	}
}

func (e *Event) Notify(s int) {
	e.sc.L.Lock()
	e.FG = len(e.fs) //重置计数器
	e.Status = s
	e.sc.L.Unlock()

	e.sc.Broadcast() //通知所有等待的订阅者

	for {
		if e.FG <= 0 {
			e.Run() //当订阅者全部运行完成后，在开启下一轮监听
			//注意这里有BUG，加入某个订阅者阻塞了，那么这里永远不会开启下一轮
			//通知方法也会卡在这个FOR语句中，最后彻底夯住。
			return
		}
	}
}

func OnCall1(e *Event, cond *sync.Cond) {
	cond.L.Lock()
	for e.FG <= 0 {
		cond.Wait() //当没有开启通知时，先暂时进入等待队列
	}
	fmt.Printf("关注香香编程喵喵喵！Status:%d\n", e.Status)
	e.FG--          //执行业务逻辑
	cond.L.Unlock() //解锁
}

func OnCall2(e *Event, cond *sync.Cond) {
	cond.L.Lock()
	for e.FG <= 0 {
		cond.Wait()
	}
	fmt.Printf("关注香香编程谢谢喵喵喵！Status:%d\n", e.Status)
	e.FG--
	cond.L.Unlock()
}

func TestObsCond() {
	l := &sync.Mutex{}
	e := &Event{
		sc: sync.NewCond(l),
		fs: []eveFun{OnCall1, OnCall2},
	}
	e.Run()

	e.Notify(1)
	e.Notify(2)
	e.Notify(3)
}

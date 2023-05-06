package code

import (
	"fmt"
	"reflect"
)

type struct1 struct {
	Name string
	Age  int
}

type struct2 struct {
	Name string
	Age  []int
}

type struct3 struct {
	Name string
	Age  *int
}

type struct4 struct {
	Name string
	Age  int
}

func TestStructs1() {
	//第一种情况 两个相同的，可比较字段的结构体
	a := struct1{"关注香香编程喵喵喵", 1}
	b := struct1{"关注香香编程喵喵喵", 1}
	fmt.Printf("第一次比较结果：%t\n", a == b) //true

	b = struct1{"关注香香编程谢谢喵喵喵", 1}
	fmt.Printf("第二次比较结果：%t\n", a == b) //false
}

func TestStructs2() {
	//第二种情况 两个相同的，不可比较字段的结构体
	//a := struct2{"关注香香编程喵喵喵", []int{1, 2, 3}}
	//b := struct2{"关注香香编程喵喵喵", []int{1, 2, 3}}
	//无法通过编译
	//fmt.Printf("第一次比较结果：%v\n", a = b)
}

func TestStructs3() {
	//第三种情况 两个相同的，带有指针字段的结构体
	a := struct3{"关注香香编程喵喵喵", new(int)}
	b := struct3{"关注香香编程喵喵喵", new(int)}
	fmt.Printf("第一次比较结果：%t\n", a == b) //false

	i := new(int)
	a = struct3{"关注香香编程喵喵喵", i}
	b = struct3{"关注香香编程喵喵喵", i}
	fmt.Printf("第一次比较结果：%t\n", a == b) //true
}

func TestStructs4() {
	//第四种情况 两个拥有相同的，可比较字段的结构体 无法通过编译
	//a := struct1{"关注香香编程喵喵喵", 1}
	//b := struct4{"关注香香编程喵喵喵", 1}
	//fmt.Printf("第一次比较结果：%t\n", a == b)
}

func TestStructs5() {
	//使用反射来判断结构体是否相同
	a := struct1{"关注香香编程喵喵喵", 1}
	b := struct1{"关注香香编程喵喵喵", 1}
	fmt.Printf("第一次比较结果：%t\n", reflect.DeepEqual(a, b)) //true

	c := struct2{"关注香香编程喵喵喵", []int{1, 2, 3}}
	d := struct2{"关注香香编程喵喵喵", []int{1, 2, 3}}
	fmt.Printf("第二次比较结果：%t\n", reflect.DeepEqual(c, d)) //true

	z := struct3{"关注香香编程喵喵喵", new(int)}
	y := struct3{"关注香香编程喵喵喵", new(int)}
	fmt.Printf("第三次比较结果：%t\n", reflect.DeepEqual(z, y)) //true

	e := struct1{"关注香香编程喵喵喵", 1}
	g := struct4{"关注香香编程喵喵喵", 1}
	fmt.Printf("第四次比较结果：%t\n", reflect.DeepEqual(e, g)) //false
}

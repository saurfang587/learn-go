package code

import "fmt"

func TestGeneric1[T string | int64](s T) T {
	fmt.Printf("关注香香编程喵喵喵！%v", s)
	return s
}

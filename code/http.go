package code

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
)

func TestHttp1() {
	url := "https://www.bilibili.com"
	resp, err := http.Get(url)
	defer func() {
		_ = resp.Body.Close()
	}()

	if err != nil {
		fmt.Printf("请求错误: %v\n", err)
	}

	fmt.Printf("返回状态:\n%s\n", resp.Status)
}

func TestHttp2() {
	url := "https://www.bilibili.com"
	req, _ := http.NewRequest(http.MethodGet, url, nil) //req 是一个Request结构，它有大量的方法的熟悉 可以自定义。
	req.Form.Add("test", "1231")                        //构造一个表单提交
	req.Header.Set("Cookie", "123")                     //设置Cookie

	resp, err := http.DefaultClient.Do(req) //这里使用的依然是默认的DefaultClient
	if err != nil {
		fmt.Printf("请求错误: %v\n", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	fmt.Printf("返回状态:\n%s\n", resp.Status)
}

func TestHttp3() {
	conn, err := net.Dial("tcp", "bilibili.com:80")
	if err != nil {
		fmt.Printf("connect err => %s\n", err.Error())
	}
	buf := bytes.Buffer{}
	buf.WriteString("GET / HTTP/1.1\r\n")
	buf.WriteString("Host: baidu.com\r\n")
	buf.WriteString("USer-Agent: Go-http-client/1.1\r\n")
	// 请求头结束
	buf.WriteString("\r\n")
	// 请求body结束
	buf.WriteString("\r\n\r\n")
	_, _ = conn.Write(buf.Bytes())
	// 获取响应信息
	resp, _ := io.ReadAll(conn)
	fmt.Printf("响应信息\n%q", resp)
}

func TestHttp4() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprintf(writer, "关注香香编程喵喵喵，关注香香编程谢谢喵喵喵！")
	})
	panic(http.ListenAndServe(":8081", nil))
}

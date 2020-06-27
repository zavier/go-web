package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func main() {
	// 生成地址
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:9999")
	checkError(err)
	log.Println("ResolveTCPAddr success")
	// 建立连接
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	defer conn.Close()
	checkError(err)
	log.Println("DialTCP success")
	// 写入
	_, err = conn.Write([]byte("HEAD / HTTP/1.0"))
	checkError(err)
	// 读取
	result, err := ioutil.ReadAll(conn)
	checkError(err)
	fmt.Println("client read data:" + string(result))
	os.Exit(0)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
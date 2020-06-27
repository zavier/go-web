package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	service := ":9999"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	// 监听端口连接
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		// 有新连接进入
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	bytes, _ := ioutil.ReadAll(conn)
	log.Println("read data:" + string(bytes))
	// 返回当前时间
	daytime := time.Now().String()
	n, err := conn.Write([]byte(daytime))
	checkError(err)
	log.Printf("server response %d bytes data", n)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
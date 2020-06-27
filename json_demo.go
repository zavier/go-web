package main

import (
	"encoding/json"
	"fmt"
)

type JsonServer struct {
	ServerName string
	ServerIP string
}

type Serverslice struct {
	Servers []JsonServer
}

func main() {
	var s Serverslice
	str := `{"servers":[{"serverName":"Shanghai","serverIP":"127.0.0.1"},{"serverName":"beijing","serverIP":"127.0.0.2"}]}`
	json.Unmarshal([]byte(str), &s)
	fmt.Println(s)
}
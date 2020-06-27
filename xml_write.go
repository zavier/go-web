package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Servers struct {
	XMLName     xml.Name `xml:"servers"`
	Version     string   `xml:"version,attr"`
	Svs         []Server `xml:"server"`
}

type Server struct {
	ServerName string   `xml:"serverName"`
	ServerIP   string   `xml:"serverIP"`
}

func writeXml() {
	v := &Servers{Version:"1"}
	v.Svs = append(v.Svs, Server{"shanghai", "127.0.0.1"})
	v.Svs = append(v.Svs, Server{"beijing", "127.0.0.2"})
	output, err := xml.MarshalIndent(v, " ", "	")
	if err != nil {
		fmt.Println("error: %v\n", err)
	}
	os.Stdout.Write([]byte(xml.Header))
	os.Stdout.Write(output)
}

func main() {
	writeXml()
}
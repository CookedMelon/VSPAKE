package main

import (
	"VSPAKE/client"
	"VSPAKE/common"
	"VSPAKE/server"
	"fmt"
)

func hex(b []byte) string {
	hex := "["
	for _, v := range b {
		hex += fmt.Sprintf(" 0x%02x ,", v)
	}
	hex += "]"
	return hex
}

func main() {
	common.LoadPasswd()
	client := client.InitClient("passwd", "client_user_Alice")
	server := server.InitClient(common.GetHashedPasswd("server_user_Bob", "client_user_Alice"), "server_user_Bob")
	recvthing := client.ClientHello()
	fmt.Printf("client hello:%#v\n", recvthing)
	server.RecvHelloMessage(recvthing)
	recvthing = server.ServerKeyExchange()
	fmt.Println("server key exchange:", hex(recvthing))
	fmt.Println(client)
	client.Update(recvthing)
	fmt.Println(client)
}

package main

import (
	"VSPAKE/check"
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
func Autopolish(s string) []byte {
	b := make([]byte, 32, 0xff)
	copy(b, s)
	return b
}
func main() {
	cliname := Autopolish("client_user_Alice")
	servname := Autopolish("server_user_Bob")
	passwd := Autopolish("password")
	hkey := common.GetHashKey(servname, cliname, passwd)
	_, client := client.InitClient(passwd, cliname)
	_, server := server.InitClient(hkey, servname)
	recvthing := client.ClientHello()
	server.RecvHelloMessage(recvthing)
	recvthing = server.ServerKeyExchange()
	client.Update(recvthing)
	recvthing = client.SendClientKeyExchange()
	server.Update(recvthing)
	server.PrintpK()
	client.PrintpK()
	fmt.Println("K eqal?", check.CheckK(server, client))

}

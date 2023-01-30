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
	client.GetAuthentKeys()
	server.GetAuthentKeys()
	fmt.Println("PreMasterSecret eqal?", check.CheckPreMasterSecret(server, client))
	fmt.Println("akey eqal?", check.Checkakey(server, client))

	//explicit mutual aythentication
	client.GetKDFs()
	server.GetKDFs()
	recvthing = server.SendKDF1()
	ans := client.AuthenticateKDF1(recvthing)
	if ans {
		fmt.Println("KDF1 pass")
	} else {
		fmt.Println("KDF1 fail")
		return
	}
	recvthing = client.SendKDF2()
	ans = server.AuthenticateKDF2(recvthing)
	if ans {
		fmt.Println("KDF2 pass")
	} else {
		fmt.Println("KDF2 fail")
		return
	}

	//Compute master secret and session key
	client.GetMasterSecretAndKey()
	server.GetMasterSecretAndKey()
	fmt.Println("MasterSecret eqal?", check.CheckMasterSecret(server, client))
	fmt.Println("SessionKey eqal?", check.CheckSessionKey(server, client))
	//send and recv message
	cAead, _ := client.Getgcm()
	sAead, _ := server.Getgcm()

	recvthing = client.SendText([]byte("hello"), cAead)
	detext, _ := server.DecryptText(recvthing, sAead)
	fmt.Println("server recv", string(detext))

	recvthing = server.SendText([]byte("world"), sAead)
	detext, _ = client.DecryptText(recvthing, cAead)
	fmt.Println("client recv", string(detext))
}

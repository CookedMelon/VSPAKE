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
	b := make([]byte, 32)
	for i := 0; i < len(b); i++ {
		b[i] = 0xff
	}
	// fmt.Println("b", b)
	copy(b, s)
	return b
}

func main() {
	// for ii := 0; ii < 100; ii++ {
	fmt.Println("-----------------------")
	cliname := Autopolish("client_user_Alice")
	servname := Autopolish("server_user_Bob")
	passwd := Autopolish("password")
	hkey := common.GetHashKey(servname, cliname, passwd)
	_, client := client.InitClient(passwd, cliname)
	_, server := server.InitClient(hkey, servname)
	recvthing := client.ClientHello()
	// fmt.Println("client hello", hex(recvthing))
	server.RecvHelloMessage(recvthing)
	recvthing = server.ServerKeyExchange()
	// fmt.Println("server key exchange", hex(recvthing))
	client.Update(recvthing)
	recvthing = client.SendClientKeyExchange()
	// fmt.Println("client key exchange", hex(recvthing))
	server.Update(recvthing)
	// server.PrintpK()
	// client.PrintpK()
	keq := check.CheckK(server, client)
	if !keq {
		server.Printall()
		client.Printall()
		return
	}
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
		// return
	}
	recvthing = client.SendKDF2()
	ans = server.AuthenticateKDF2(recvthing)
	if ans {
		fmt.Println("KDF2 pass")
	} else {
		fmt.Println("KDF2 fail")
		// return
	}

	//Compute master secret and session key
	client.GetMasterSecretAndKey()
	server.GetMasterSecretAndKey()
	fmt.Println("MasterSecret eqal?", check.CheckMasterSecret(server, client))
	fmt.Println("SessionKey eqal?", check.CheckSessionKey(server, client))

	cz := client.Checkzero()
	if !cz {
		fmt.Println("client check zero fail")
		// return
	}
	cz = server.Checkzero()
	if !cz {
		fmt.Println("client check zero fail")
		// return
	}
	//send and recv message
	cAead, _ := client.Getgcm()
	sAead, _ := server.Getgcm()

	recvthing = client.SendText([]byte("hello"), cAead)
	detext, _ := server.DecryptText(recvthing, sAead)
	fmt.Println("server recv", string(detext))

	recvthing = server.SendText([]byte("world"), sAead)
	detext, _ = client.DecryptText(recvthing, cAead)
	fmt.Println("client recv", string(detext))
	// }
}

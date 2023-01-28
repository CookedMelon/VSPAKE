package test

import (
	"VSPAKE/common"
	"VSPAKE/packages/elliptic"
	"crypto/rand"
	"fmt"
	"testing"
)

func TestH(t *testing.T) {
	string1 := "server_user_Bob" //
	serverbytes := make([]byte, 32)
	copy(serverbytes, string1)
	fmt.Printf("server_name:%#v\n", serverbytes)
	string2 := "client_user_Alice" //
	clientbytes := make([]byte, 32)
	copy(clientbytes, string2)
	fmt.Printf("client_name:%#v\n", clientbytes)
	string3 := "password" //
	keybytes := make([]byte, 32)
	copy(keybytes, string3)
	hkey := make([]byte, 32)
	copy(hkey, common.GetHashKey(clientbytes, serverbytes, keybytes))
	fmt.Printf("hkey       :%#v\n", hkey)
}
func TestGetRandPoint(t *testing.T) {
	r := make([]byte, 32)
	rand.Read(r)
	curve := new(elliptic.CurveDetail)
	curve.Init()
	p := curve.Mult(&curve.BasePoint, r)
	fmt.Printf("p:%#v\n", p)
}
func TestCheckOn(t *testing.T) {
	curve := new(elliptic.CurveDetail)
	curve.Init()
	fmt.Println(curve.IfOnCurve(curve.P))
	fmt.Println(curve.IfOnCurve(curve.Q))
}

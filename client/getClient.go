package client

import (
	"math/big"
	"unsafe"
)

type CiphSel struct {
	// contains filtered or unexported fields
	p, P, Q *big.Int //p,q为素数，P,Q为椭圆曲线上的点

}
type Client struct {
	// contains filtered or unexported fields
	//public
	G        CiphSel
	nc, name []byte
	//private
	key []byte
}
type HelloMessage struct {
	nc, name []byte
}

func GetClient(key, name, nc []byte) Client {
	return Client{key: key, name: name, nc: nc}
}
func GetClientHelloMessage(client *Client, hm *HelloMessage) {
	hm.nc = client.nc
	hm.name = client.name
}
func ClientHello(client *Client) {
	message := &HelloMessage{}
	GetClientHelloMessage(client, message)
	Send(*(*[]byte)(unsafe.Pointer(message)))
}

func Send(b []byte) []byte {
	return b
}

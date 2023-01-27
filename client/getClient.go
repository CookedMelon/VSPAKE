package client

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/json"
	"errors"
	"server"
	"unsafe"
)

func GetClientHelloMessage(client *Client) *HelloMessage {
	hm := &HelloMessage{}
	hm.nc = client.nc
	hm.name = client.name
	return hm
}
func GetClient(key, name, nc []byte) *Client {
	return &Client{key: key, name: name, nc: nc}
}

func (client *Client) ClientHello() {
	message := GetClientHelloMessage(client)
	Send(*(*[]byte)(unsafe.Pointer(message)))
}

func (client *Client) Update(climsg []byte) error {
	kec := server.KeyExchangeMsg{}
	err := json.Unmarshal(climsg, &kec)
	if err != nil {
		return errors.New("json unmarshal error in client update/1")
	}
	err = json.Unmarshal(kec.Rbyte, client.R)
	if err != nil {
		return errors.New("json unmarshal error in client update/2")
	}
	err = json.Unmarshal(kec.Ybyte, client.Y)
	if err != nil {
		return errors.New("json unmarshal error in client update/3")
	}
	hasher := md5.New()
	hasher.Write(client.cname)
	hasher.Write(client.sname)
	hasher.Write(client.key)
	H0 := hasher.Sum(nil)
	client.x = make([]byte, 32)
	rand.Read(client.x)
	temp1 := client.Mult(client.G.P, kec.X)
	temp2 := client.Mult(client.G.R, H0)
	client.K = client.Add(temp1, temp2)

	return nil
}

func Send(b []byte) []byte {
	return b
}

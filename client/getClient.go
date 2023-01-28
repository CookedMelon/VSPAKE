package client

import (
	"VSPAKE/server"
	"crypto/md5"
	"crypto/rand"
	"encoding/json"
	"errors"
	"unsafe"
)

func GetClientHelloMessage(client *Client) *HelloMessage {
	hm := &HelloMessage{}
	hm.nc = client.NC
	hm.name = client.Cname
	return hm
}
func InitClient(key, cname []byte) *Client {
	nc := make([]byte, 32)
	rand.Read(nc)
	return &Client{key: key, Cname: cname, NC: nc}
}

func (client *Client) ClientHello() {
	message := GetClientHelloMessage(client)
	Send(*(*[]byte)(unsafe.Pointer(message)))
}
func (client *Client) RecvHelloMessage(mess []byte) {
	client.NC = mess
}
func (client *Client) Update(climsg []byte) error {
	kec := server.KeyExchangeMsg{}
	err := json.Unmarshal(climsg, &kec)
	if err != nil {
		return errors.New("json unmarshal error in client update/1")
	}
	err = json.Unmarshal(kec.Rbyte, client.pR)
	if err != nil {
		return errors.New("json unmarshal error in client update/2")
	}
	err = json.Unmarshal(kec.Ybyte, client.pY)
	if err != nil {
		return errors.New("json unmarshal error in client update/3")
	}
	hasher := md5.New()
	hasher.Write(client.Cname)
	hasher.Write(client.Sname)
	hasher.Write(client.key)
	H0 := hasher.Sum(nil)
	client.x = make([]byte, 32)
	rand.Read(client.x)
	temp1 := client.G.Mult(client.G.P, client.x)
	temp2 := client.G.Mult(client.pR, H0)
	client.pX = client.G.Add(temp1, temp2)
	client.pK = client.G.Mult(client.pY, client.x)

	return nil
}
func (client *Client) GetAuthentKeys() {
	hasher := md5.New()
	hasher.Write(client.Cname)
	hasher.Write(client.Sname)
	hasher.Write(client.NC)
	hasher.Write(client.NS)
	hasher.Write(*(*[]byte)(unsafe.Pointer(client.pX)))
	hasher.Write(*(*[]byte)(unsafe.Pointer(client.pR)))
	hasher.Write(*(*[]byte)(unsafe.Pointer(client.pY)))
	hasher.Write(*(*[]byte)(unsafe.Pointer(client.pK)))
	client.preMasterSecret = hasher.Sum(nil)
	hasher.Reset()
	hasher.Write(client.preMasterSecret)
	hasher.Write(client.NC)
	hasher.Write(client.NS)
	client.aKey = hasher.Sum(nil)
}
func (client *Client) GetKDFs() {
	hasher := md5.New()
	hasher.Write(client.aKey)
	hasher.Write(client.Cname)
	hasher.Write(client.Sname)
	hasher.Write(client.NC)
	hasher.Write(client.NS)
	client.kdf1 = hasher.Sum(nil)
	hasher.Reset()
	hasher.Write(client.aKey)
	hasher.Write(client.Sname)
	hasher.Write(client.Cname)
	hasher.Write(client.NS)
	hasher.Write(client.NC)
	client.kdf2 = hasher.Sum(nil)
}
func (client *Client) AuthenticateKDF1(recvkdf []byte) {
	if string(client.kdf1) == string(recvkdf) {
		return
	}
	panic("kdf1 not equal")
}

func (client *Client) SendKDF2(recvkdf []byte) {
	Send(client.kdf2)
}
func (client *Client) GetMaeterSecret() {
	hasher := md5.New()
	hasher.Write(client.preMasterSecret)
	hasher.Write(client.NC)
	hasher.Write(client.NS)
	client.masterSecret = hasher.Sum(nil)
}
func (client *Client) GetSessionKey() {
	hasher := md5.New()
	hasher.Write(client.masterSecret)
	hasher.Write(client.Cname)
	hasher.Write(client.Sname)
	hasher.Write(client.NC)
	hasher.Write(client.NS)
	client.sessionKey = hasher.Sum(nil)
}
func Send(b []byte) []byte {
	return b
}

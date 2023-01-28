package client

import (
	"VSPAKE/common"
	"VSPAKE/packages/elliptic"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"math/big"
	"unsafe"
)

func InitClient(passwd, cname string) *Client {
	G := new(elliptic.CurveDetail)
	G.Init()
	fmt.Println("G:", G)
	nc := make([]byte, 32)
	key := make([]byte, 32)
	Cname := make([]byte, 32)
	copy(key, passwd)
	copy(Cname, cname)
	rand.Read(nc)
	return &Client{G: G, key: key, Cname: Cname, NC: nc}
}
func (client *Client) GetClientHelloMessage() *common.ClientHelloMessage {
	hm := &common.ClientHelloMessage{}
	copy(hm.NC[:], client.NC)
	copy(hm.Name[:], client.Cname)
	return hm
}
func (client *Client) ClientHello() []byte {
	message := client.GetClientHelloMessage()
	fmt.Println("client hello message", message)
	return common.Sendc(message)
}
func (client *Client) RecvHelloMessage(mess []byte) {
	client.NC = mess
}
func (client *Client) Update(sermsg []byte) error {
	kec := common.ServerKeyExchangeMsg{}
	copy(kec.NS[:], sermsg[:32])
	copy(kec.Sname[:], sermsg[32:64])
	copy(kec.Rbyte[:], sermsg[64:128])
	copy(kec.Ybyte[:], sermsg[128:192])
	client.NS = kec.NS[:]
	client.Sname = kec.Sname[:]
	client.pR = new(elliptic.CurvePoint)
	client.pY = new(elliptic.CurvePoint)
	client.pR.X = new(big.Int).SetBytes(kec.Rbyte[:28])
	client.pR.Y = new(big.Int).SetBytes(kec.Rbyte[32 : 32+28])
	fmt.Println("pR.X:", client.pR.X)
	fmt.Println("pR.Y:", client.pR.Y)
	client.pY.X = new(big.Int).SetBytes(kec.Ybyte[:28])
	fmt.Println("pY.X:", client.pY.X)
	client.pY.Y = new(big.Int).SetBytes(kec.Ybyte[32 : 32+28])
	client.hkey = common.GetHashKey(client.Cname, client.Sname, client.key)
	client.x = make([]byte, 32)
	rand.Read(client.x)
	temp1 := client.G.Mult(client.G.P, client.x)
	temp2 := client.G.Mult(client.pR, client.hkey)
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
	common.Send(client.kdf2)
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

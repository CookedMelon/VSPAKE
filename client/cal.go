package client

import (
	"VSPAKE/common"
	"VSPAKE/packages/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

func InitClient(passwd, cname []byte) (error, *Client) {
	if len(passwd) != 32 {
		return fmt.Errorf("passwd length error"), nil
	}
	if len(cname) != 32 {
		return fmt.Errorf("cname length error"), nil
	}
	G := new(elliptic.CurveDetail)
	G.Init()
	nc := make([]byte, 32)
	key := make([]byte, 32)
	Cname := make([]byte, 32)
	copy(key, passwd)
	copy(Cname, cname)
	rand.Read(nc)
	return nil, &Client{G: G, key: key, Cname: Cname, NC: nc}
}
func (client *Client) getclihello() *common.ClientHelloMessage {
	hm := &common.ClientHelloMessage{}
	copy(hm.NC[:], client.NC)
	copy(hm.Name[:], client.Cname)
	return hm
}
func (client *Client) ClientHello() []byte {
	message := client.getclihello()
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
	client.pY.X = new(big.Int).SetBytes(kec.Ybyte[:28])
	client.pY.Y = new(big.Int).SetBytes(kec.Ybyte[32 : 32+28])
	client.hkey = common.GetHashKey(client.Sname, client.Cname, client.key)
	client.x = make([]byte, 32)
	rand.Read(client.x)
	temp1 := client.G.Mult(client.G.P, client.x)
	temp2 := client.G.Mult(client.pR, client.hkey)
	client.pX = client.G.Add(temp1, temp2)
	client.pK = client.G.Mult(client.pY, client.x)
	return nil
}
func (client *Client) SendClientKeyExchange() []byte {

	return common.SendExcCli(client.getkeyexc())
}
func (client *Client) getkeyexc() *common.ClientKeyExchangeMsg {
	kec := common.ClientKeyExchangeMsg{}
	copy(kec.Xbyte[:], client.pX.X.Bytes())
	copy(kec.Xbyte[32:], client.pX.Y.Bytes())
	return &kec
}
func (client *Client) GetAuthentKeys() {
	hasher := sha256.New()
	hasher.Write(client.Cname)
	hasher.Write(client.Sname)
	hasher.Write(client.NC)
	hasher.Write(client.NS)
	hasher.Write(client.pX.X.Bytes())
	hasher.Write(client.pX.Y.Bytes())
	hasher.Write(client.pR.X.Bytes())
	hasher.Write(client.pR.Y.Bytes())
	hasher.Write(client.pY.X.Bytes())
	hasher.Write(client.pY.Y.Bytes())
	hasher.Write(client.pK.X.Bytes())
	hasher.Write(client.pK.Y.Bytes())

	// fmt.Println("client")
	// fmt.Println(client.Cname)
	// fmt.Println(client.Sname)
	// fmt.Println(client.NC)
	// fmt.Println(client.NS)
	// fmt.Println(client.pX.X.Bytes())
	// fmt.Println(client.pX.Y.Bytes())
	// fmt.Println(client.pR.X.Bytes())
	// fmt.Println(client.pR.Y.Bytes())
	// fmt.Println(client.pY.X.Bytes())
	// fmt.Println(client.pY.Y.Bytes())
	// fmt.Println(client.pK.X.Bytes())
	// fmt.Println(client.pK.Y.Bytes())
	// fmt.Println()

	client.preMasterSecret = hasher.Sum(nil)
	hasher.Reset()
	hasher.Write(client.preMasterSecret)
	hasher.Write(client.NC)
	hasher.Write(client.NS)
	client.aKey = hasher.Sum(nil)
}
func (client *Client) GetKDFs() {
	hasher := sha256.New()
	hasher.Write(client.aKey)
	hasher.Write(client.Cname)
	hasher.Write(client.Sname)
	hasher.Write(client.NC)
	hasher.Write(client.NS)
	client.kdf1 = hasher.Sum(nil)
	// fmt.Println("client.kdf1", client.kdf1)
	hasher.Reset()
	hasher.Write(client.aKey)
	hasher.Write(client.Sname)
	hasher.Write(client.Cname)
	hasher.Write(client.NS)
	hasher.Write(client.NC)
	client.kdf2 = hasher.Sum(nil)
}

func (client *Client) AuthenticateKDF1(recvkdf []byte) bool {
	return string(client.kdf1) == string(recvkdf)
}

func (client *Client) SendKDF2() []byte {
	return common.Send(client.kdf2)
}
func (client *Client) GetMaeterSecret() {
	hasher := sha256.New()
	hasher.Write(client.preMasterSecret)
	hasher.Write(client.NC)
	hasher.Write(client.NS)
	client.masterSecret = hasher.Sum(nil)
}
func (client *Client) GetSessionKey() {
	hasher := sha256.New()
	hasher.Write(client.masterSecret)
	hasher.Write(client.Cname)
	hasher.Write(client.Sname)
	hasher.Write(client.NC)
	hasher.Write(client.NS)
	client.sessionKey = hasher.Sum(nil)
}

func (client *Client) GetMasterSecretAndKey() {
	hasher := sha256.New()
	hasher.Write(client.preMasterSecret)
	hasher.Write(client.NC)
	hasher.Write(client.NS)
	client.masterSecret = hasher.Sum(nil)
	hasher.Reset()
	hasher.Write(client.masterSecret)
	hasher.Write(client.Cname)
	hasher.Write(client.Sname)
	hasher.Write(client.NC)
	hasher.Write(client.NS)
	client.sessionKey = hasher.Sum(nil)
}

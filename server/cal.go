package server

import (
	"VSPAKE/common"
	"VSPAKE/packages/elliptic"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

func InitClient(hkey, sname []byte) (error, *Server) {
	if len(hkey) != 32 {
		return fmt.Errorf("hkey length error"), nil
	}
	if len(sname) != 32 {
		return fmt.Errorf("sname length error"), nil
	}
	G := new(elliptic.CurveDetail)
	G.Init()
	NS := make([]byte, 32)
	Sname := make([]byte, 32)
	Hkey := make([]byte, 32)
	copy(Sname, sname)
	copy(Hkey, hkey)
	rand.Read(NS)
	V := G.Mult(G.Q, Hkey)
	return nil, &Server{G: G, Sname: Sname, NS: NS, hkey: Hkey, pV: V}
}
func (server *Server) RecvHelloMessage(climsg []byte) error {
	server.NC = make([]byte, 32)
	server.Cname = make([]byte, 32)
	copy(server.NC[:], climsg[0:32])
	copy(server.Cname[:], climsg[32:])
	return nil
}

func (server *Server) ServerKeyExchange() []byte {
	message := server.getserkeyexc()
	return common.SendExcSer(message)
}
func (server *Server) getserkeyexc() *common.ServerKeyExchangeMsg {
	server.r = make([]byte, 32)
	server.y = make([]byte, 32)
	rand.Read(server.r)
	rand.Read(server.y)
	mess := &common.ServerKeyExchangeMsg{}
	copy(mess.NS[:], server.NS)
	copy(mess.Sname[:], server.Sname)
	server.pR = server.GetR()
	server.pY = server.GetY()
	copy(mess.Rbyte[:], server.pR.X.Bytes())
	copy(mess.Rbyte[32:], server.pR.Y.Bytes())
	copy(mess.Ybyte[:], server.pY.X.Bytes())
	copy(mess.Ybyte[32:], server.pY.Y.Bytes())
	return mess
}
func (server *Server) GetR() *elliptic.CurvePoint {
	return server.G.Mult(server.G.Q, server.r)
}
func (server *Server) GetY() *elliptic.CurvePoint {
	return server.G.Mult(server.G.P, server.y)
}
func (server *Server) Update(Xbyte []byte) {
	X := elliptic.GetEmptyCurvePoint()
	X.X.SetBytes(Xbyte[0:28])
	X.Y.SetBytes(Xbyte[32 : 32+28])
	server.pX = X
	temp := server.G.Mult(server.pV, server.r)
	server.G.GetNeg(temp)
	temp2 := server.G.Add(server.pX, temp)
	server.pK = server.G.Mult(temp2, server.y)
}
func (server *Server) GetAuthentKeys() {
	hasher := sha256.New()
	hasher.Write(server.Cname)
	hasher.Write(server.Sname)
	hasher.Write(server.NC)
	hasher.Write(server.NS)
	hasher.Write(server.pX.X.Bytes())
	hasher.Write(server.pX.Y.Bytes())
	hasher.Write(server.pR.X.Bytes())
	hasher.Write(server.pR.Y.Bytes())
	hasher.Write(server.pY.X.Bytes())
	hasher.Write(server.pY.Y.Bytes())
	hasher.Write(server.pK.X.Bytes())
	hasher.Write(server.pK.Y.Bytes())
	// fmt.Println("server")
	// fmt.Println(server.Cname)
	// fmt.Println(server.Sname)
	// fmt.Println(server.NC)
	// fmt.Println(server.NS)
	// fmt.Println(server.pX.X.Bytes())
	// fmt.Println(server.pX.Y.Bytes())
	// fmt.Println(server.pR.X.Bytes())
	// fmt.Println(server.pR.Y.Bytes())
	// fmt.Println(server.pY.X.Bytes())
	// fmt.Println(server.pY.Y.Bytes())
	// fmt.Println(server.pK.X.Bytes())
	// fmt.Println(server.pK.Y.Bytes())
	// fmt.Println()

	server.preMasterSecret = hasher.Sum(nil)
	hasher.Reset()
	hasher.Write(server.preMasterSecret)
	hasher.Write(server.NC)
	hasher.Write(server.NS)
	server.aKey = hasher.Sum(nil)
}
func (server *Server) GetKDFs() {
	hasher := sha256.New()
	hasher.Write(server.aKey)
	hasher.Write(server.Cname)
	hasher.Write(server.Sname)
	hasher.Write(server.NC)
	hasher.Write(server.NS)
	server.kdf1 = hasher.Sum(nil)
	// fmt.Println("server.kdf1", server.kdf1)
	hasher.Reset()
	hasher.Write(server.aKey)
	hasher.Write(server.Sname)
	hasher.Write(server.Cname)
	hasher.Write(server.NS)
	hasher.Write(server.NC)
	server.kdf2 = hasher.Sum(nil)
}
func (server *Server) SendKDF1() []byte {
	return common.Send(server.kdf1)
}
func (server *Server) AuthenticateKDF2(recvkdf []byte) bool {
	return string(server.kdf2) == string(recvkdf)
}

func (server *Server) GetMasterSecretAndKey() {
	hasher := sha256.New()
	hasher.Write(server.preMasterSecret)
	hasher.Write(server.NC)
	hasher.Write(server.NS)
	server.masterSecret = hasher.Sum(nil)
	hasher.Reset()
	hasher.Write(server.masterSecret)
	hasher.Write(server.Cname)
	hasher.Write(server.Sname)
	hasher.Write(server.NC)
	hasher.Write(server.NS)
	server.sessionKey = hasher.Sum(nil)

}
func (server *Server) Getgcm() (cipher.AEAD, error) {
	block, err := aes.NewCipher(server.sessionKey)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return aead, nil
}
func (server *Server) SendText(plaintext []byte, aead cipher.AEAD) []byte {
	nonce := make([]byte, 12)
	rand.Read(nonce)
	return append(nonce, aead.Seal(nil, nonce, plaintext, server.aKey)...)
}
func (server *Server) DecryptText(ciphertext []byte, aead cipher.AEAD) ([]byte, error) {
	nonce := ciphertext[:12]
	return aead.Open(nil, nonce, ciphertext[12:], server.aKey)
}

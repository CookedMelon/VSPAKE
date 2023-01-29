package server

import (
	"VSPAKE/common"
	"VSPAKE/packages/elliptic"
	"crypto/rand"
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

package server

import (
	"VSPAKE/common"
	"VSPAKE/packages/elliptic"
	"crypto/rand"
	"fmt"
)

func InitClient(hkey, sname string) *Server {
	G := new(elliptic.CurveDetail)
	G.Init()

	NS := make([]byte, 32)
	Sname := make([]byte, 32)
	Hkey := make([]byte, 32)
	copy(Sname, sname)
	copy(Hkey, hkey)
	rand.Read(NS)
	return &Server{G: G, Sname: Sname, NS: NS, hkey: Hkey}
}
func (server *Server) RecvHelloMessage(climsg []byte) error {
	copy(server.NC[:], climsg[0:32])
	copy(server.Cname[:], climsg[32:])
	return nil
}

func (server *Server) ServerKeyExchange() []byte {
	message := server.GetServerKeyExchangeMessage()
	return common.SendExc(message)
}
func (server *Server) GetServerKeyExchangeMessage() *common.ServerKeyExchangeMsg {
	server.r = make([]byte, 32)
	server.y = make([]byte, 32)
	rand.Read(server.r)
	rand.Read(server.y)
	mess := &common.ServerKeyExchangeMsg{}
	copy(mess.NS[:], server.NS)
	copy(mess.Sname[:], server.Sname)
	server.pR = server.GetR()
	fmt.Println("pR.X-", server.pR.X)
	fmt.Println("pR.Y-", server.pR.Y)
	server.pY = server.GetY()
	fmt.Println(server.pR.X.Bytes(), server.pR.Y.Bytes())
	fmt.Println(server.pY.X.Bytes(), server.pY.Y.Bytes())
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

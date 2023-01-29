package server

import "fmt"

func (server *Server) PrintpK() {
	fmt.Println("server pK.X-", server.pK.X)
	fmt.Println("server pK.Y-", server.pK.Y)
}

func (server *Server) Printall() {
	fmt.Println("server pX-", server.pX)
	fmt.Println("server pR-", server.pR)
	fmt.Println("server pY-", server.pY)
	fmt.Println("server pK-", server.pK)
	fmt.Println("server hk-", server.hkey)
}

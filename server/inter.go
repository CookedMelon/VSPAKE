package server

func (server *Server) OutputK() []byte {
	k := make([]byte, 64)
	copy(k[0:32], server.pK.X.Bytes())
	copy(k[32:64], server.pK.Y.Bytes())
	return k
}

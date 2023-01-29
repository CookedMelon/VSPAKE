package server

func (server *Server) OutputK() []byte {
	k := make([]byte, 64)
	copy(k[0:32], server.pK.X.Bytes())
	copy(k[32:64], server.pK.Y.Bytes())
	return k
}
func (server *Server) OutputPreMasterSecret() []byte {
	s := make([]byte, 32)
	copy(s[0:32], server.preMasterSecret)
	return s
}
func (server *Server) Outputakey() []byte {
	x := make([]byte, 32)
	copy(x[0:32], server.aKey)
	return x
}
func (server *Server) OutputMasterSecret() []byte {
	x := make([]byte, 32)
	copy(x[0:32], server.masterSecret)
	return x
}
func (server *Server) OutputSessionKey() []byte {
	x := make([]byte, 32)
	copy(x[0:32], server.sessionKey)
	return x
}

package common

type ServerKeyExchangeMsg struct {
	NS, Sname    [32]byte
	Rbyte, Ybyte [64]byte
}
type ClientKeyExchangeMsg struct {
	Xbyte [64]byte
}
type ClientHelloMessage struct {
	NC, Name [32]byte
}
type ServerHelloMessage struct {
	NS [32]byte
}

func Send(b []byte) []byte {
	return b
}
func SendExcSer(msg *ServerKeyExchangeMsg) []byte {
	ret := make([]byte, 32*2+64*2)
	copy(ret[0:32], msg.NS[:])
	copy(ret[32:64], msg.Sname[:])
	copy(ret[64:128], msg.Rbyte[:])
	copy(ret[128:192], msg.Ybyte[:])
	return ret[:]
}
func SendExcCli(msg *ClientKeyExchangeMsg) []byte {
	ret := make([]byte, 64)
	copy(ret[0:64], msg.Xbyte[:])
	return ret[:]
}
func SendHellos(msg *ServerHelloMessage) []byte {
	ret := make([]byte, 32)
	copy(ret[0:32], msg.NS[:])
	return ret[:]
}
func Sendc(msg *ClientHelloMessage) []byte {
	ret := make([]byte, 32*2)
	copy(ret[0:32], msg.NC[:])
	copy(ret[32:64], msg.Name[:])
	return ret[:]
}
func CopyPoingByte(dst, src []byte) {
	copy(dst[28-len(src):], src)
}

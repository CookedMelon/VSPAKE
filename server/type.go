package server

type KeyExchangeMsg struct {
	ns, s        []byte
	Rbyte, Ybyte []byte
}

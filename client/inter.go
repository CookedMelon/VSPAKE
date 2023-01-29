package client

func (client *Client) OutputK() []byte {
	k := make([]byte, 64)
	copy(k[0:32], client.pK.X.Bytes())
	copy(k[32:64], client.pK.Y.Bytes())
	return k
}
func (client *Client) OutputPreMasterSecret() []byte {
	s := make([]byte, 32)
	copy(s[0:32], client.preMasterSecret)
	return s
}
func (client *Client) Outputakey() []byte {
	s := make([]byte, 32)
	copy(s[0:32], client.aKey)
	return s
}

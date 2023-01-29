package client

func (client *Client) OutputK() []byte {
	k := make([]byte, 64)
	copy(k[0:32], client.pK.X.Bytes())
	copy(k[32:64], client.pK.Y.Bytes())
	return k
}

package client

import "fmt"

func (client *Client) PrintpK() {
	fmt.Println("client pK.X-", client.pK.X)
	fmt.Println("client pK.Y-", client.pK.Y)
}

func (client *Client) Printall() {
	fmt.Println("client pX-", client.pX)
	fmt.Println("client pR-", client.pR)
	fmt.Println("client pY-", client.pY)
	fmt.Println("client pK-", client.pK)
	fmt.Println("client hk-", client.hkey)
}
func (client *Client) OutputMasterSecret() []byte {
	x := make([]byte, 32)
	copy(x[0:32], client.masterSecret)
	return x
}
func (client *Client) OutputSessionKey() []byte {
	x := make([]byte, 32)
	copy(x[0:32], client.sessionKey)
	return x
}

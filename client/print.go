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

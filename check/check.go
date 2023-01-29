package check

import (
	"VSPAKE/client"
	"VSPAKE/server"
	"fmt"
)

func CheckK(s *server.Server, c *client.Client) bool {
	sk := s.OutputK()
	ck := c.OutputK()
	for i := 0; i < 64; i++ {
		if sk[i] != ck[i] {
			return false
		}
	}
	return true
}
func CheckPreMasterSecret(s *server.Server, c *client.Client) bool {
	sk := s.OutputPreMasterSecret()
	ck := c.OutputPreMasterSecret()
	fmt.Println("sk", sk)
	fmt.Println("ck", ck)
	for i := 0; i < 32; i++ {
		if sk[i] != ck[i] {
			return false
		}
	}
	return true
}
func Checkakey(s *server.Server, c *client.Client) bool {
	sk := s.Outputakey()
	ck := c.Outputakey()
	for i := 0; i < 32; i++ {
		if sk[i] != ck[i] {
			return false
		}
	}
	return true
}
func CheckMasterSecret(s *server.Server, c *client.Client) bool {
	sk := s.OutputMasterSecret()
	ck := c.OutputMasterSecret()
	for i := 0; i < 32; i++ {
		if sk[i] != ck[i] {
			return false
		}
	}
	return true
}
func CheckSessionKey(s *server.Server, c *client.Client) bool {
	sk := s.OutputSessionKey()
	ck := c.OutputSessionKey()
	fmt.Println("sk SessionKey", sk)
	fmt.Println("ck SessionKey", ck)
	for i := 0; i < 32; i++ {
		if sk[i] != ck[i] {
			return false
		}
	}
	return true
}

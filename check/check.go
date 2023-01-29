package check

import (
	"VSPAKE/client"
	"VSPAKE/server"
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

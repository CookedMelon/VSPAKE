package client

import (
	"VSPAKE/packages/elliptic"
)

type Client struct {
	// contains filtered or unexported fields
	//public
	G                    *elliptic.CurveDetail
	NC, NS, Sname, Cname []byte
	//private
	key                           []byte //短密码
	x                             []byte //随机数
	pX, pY, pR, pK                *elliptic.CurvePoint
	sessionKey                    []byte
	preMasterSecret, masterSecret []byte
	aKey, skey                    []byte
	kdf1, kdf2                    []byte
}
type HelloMessage struct {
	nc, name []byte
}

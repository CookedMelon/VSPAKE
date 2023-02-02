package client

import (
	"VSPAKE/packages/elliptic"
)

type Client struct {
	// contains filtered or unexported fields
	//public
	G                    *elliptic.CurveDetail
	NC, NS, Sname, Cname []byte //NC NS随机
	//private
	key, hkey                     []byte //客户独有的短密码和服务端存储的哈希后的密码
	x                             []byte //随机数
	pX, pY, pR, pK                *elliptic.CurvePoint
	sessionKey                    []byte
	preMasterSecret, masterSecret []byte
	aKey                          []byte
	kdf1, kdf2                    []byte
}

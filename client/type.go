package client

import (
	"math/big"
	"packages/elliptic"
)

type CiphSel struct {
	// contains filtered or unexported fields
	P, Q *elliptic.CurvePoint //P,Q为椭圆曲线上的点
	p    *big.Int             //p为素数

}
type Client struct {
	// contains filtered or unexported fields
	//public
	G                    *CiphSel
	nc, ns, sname, cname []byte
	//private
	key             []byte   //短密码
	x               *big.Int //随机数
	X, Y, R, K      *elliptic.CurvePoint
	sessionKey      []byte
	PreMasterSecret []byte
	AKey            []byte
}
type HelloMessage struct {
	nc, name []byte
}

package server

import "VSPAKE/packages/elliptic"

type Server struct {
	// contains filtered or unexported fields
	//public
	G                    *elliptic.CurveDetail
	NC, NS, Sname, Cname []byte
	//private
	hkey                          []byte //hash后的短密码
	r, y                          []byte //随机数
	pX, pY, pR, pK, pV            *elliptic.CurvePoint
	sessionKey                    []byte
	preMasterSecret, masterSecret []byte
	aKey, skey                    []byte
	kdf1, kdf2                    []byte
}

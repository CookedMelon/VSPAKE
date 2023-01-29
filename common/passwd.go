package common

import (
	"crypto/sha256"
)

func GetHashKey(Cname, Sname, key []byte) []byte {
	hkey := make([]byte, 32, 0xff)

	hasher := sha256.New()
	hasher.Write(Cname)
	hasher.Write(Sname)
	hasher.Write(key)
	copy(hkey, hasher.Sum(nil))
	return hkey
}

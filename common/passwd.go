package common

import "crypto/md5"

func GetHashKey(Cname, Sname, key []byte) []byte {
	hkey := make([]byte, 32, 0xff)

	hasher := md5.New()
	hasher.Write(Cname)
	hasher.Write(Sname)
	hasher.Write(key)
	copy(hkey, hasher.Sum(nil))
	return hkey
}

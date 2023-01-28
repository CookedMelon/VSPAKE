package common

import "crypto/md5"

var passwd map[string]map[string]string

func LoadPasswd() {
	passwd = make(map[string]map[string]string)
	passwd["server_user_Bob"] = make(map[string]string)
	passwd["server_user_Bob"]["client_user_Alice"] = string(GetHashKey([]byte("server_user_Bob"), []byte("client_user_Alice"), []byte("password")))
}
func GetHashedPasswd(Cname, Sname string) string {
	return passwd[Sname][Cname]
}
func GetHashKey(Cname, Sname, key []byte) []byte {
	hkey := make([]byte, 32)
	hasher := md5.New()
	hasher.Write(Cname)
	hasher.Write(Sname)
	hasher.Write(key)
	copy(hkey, hasher.Sum(nil))
	return hkey
}

package test

import (
	"VSPAKE/common"
	"VSPAKE/packages/elliptic"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"testing"
)

func TestH(t *testing.T) {
	string1 := "server_user_Bob" //
	serverbytes := make([]byte, 32)
	copy(serverbytes, string1)
	fmt.Printf("server_name:%#v\n", serverbytes)
	string2 := "client_user_Alice" //
	clientbytes := make([]byte, 32)
	copy(clientbytes, string2)
	fmt.Printf("client_name:%#v\n", clientbytes)
	string3 := "password" //
	keybytes := make([]byte, 32)
	copy(keybytes, string3)
	hkey := make([]byte, 32)
	copy(hkey, common.GetHashKey(clientbytes, serverbytes, keybytes))
	fmt.Printf("hkey       :%#v\n", hkey)
}
func TestGetRandPoint(t *testing.T) {
	r := make([]byte, 32)
	rand.Read(r)
	curve := new(elliptic.CurveDetail)
	curve.Init()
	p := curve.Mult(&curve.BasePoint, r)
	fmt.Printf("p:%#v\n", p)
}
func TestCheckOn(t *testing.T) {
	curve := new(elliptic.CurveDetail)
	curve.Init()
	fmt.Println(curve.IfOnCurve(curve.P))
	fmt.Println(curve.IfOnCurve(curve.Q))
}

type AEAD interface {
	// 返回提供给Seal和Open方法的随机数nonce的字节长度
	NonceSize() int
	// 返回原始文本和加密文本的最大长度差异
	Overhead() int
	// 加密并认证明文，认证附加的data，将结果添加到dst，返回更新后的切片。
	// nonce的长度必须是NonceSize()字节，且对给定的key和时间都是独一无二的。
	// plaintext和dst可以是同一个切片，也可以不同。
	Seal(dst, nonce, plaintext, data []byte) []byte
	// 解密密文并认证，认证附加的data，如果认证成功，将明文添加到dst，返回更新后的切片。
	// nonce的长度必须是NonceSize()字节，nonce和data都必须和加密时使用的相同。
	// ciphertext和dst可以是同一个切片，也可以不同。
	Open(dst, nonce, ciphertext, data []byte) ([]byte, error)
}

func TestGCM(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	plaintext := []byte("hello world")
	nonce := []byte("123456789012")
	aead, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(aead)
	if err != nil {
		panic(err.Error())
	}
	// add := []byte("12345678901234567890123456789012")
	// add2 := []byte("12345678902234567890123456789012")
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	fmt.Printf("%s\n", ciphertext)
	ciphertext[5] = 0x44
	plaintext, err = gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%s\n", plaintext)
}

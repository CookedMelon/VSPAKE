package aes

func XORenbuf(enbuf []byte, len int) []byte {
	ans := make([]byte, 16)
	for i := 0; i < len; i++ {
		ans[i] ^= enbuf[i%16]
	}
	return ans
}

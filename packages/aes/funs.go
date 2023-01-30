package aes

type AES_KEY struct {
	RoundKey [16 * 11]byte
	Iv       [16]byte
}

func SubBytes(state []byte) {
	for i := 0; i < 16; i++ {
		state[i] = stb[state[i]]
	}
}
func InvSubBytes(state []byte) {
	for i := 0; i < 16; i++ {
		state[i] = stdinv[state[i]]
	}
}
func ShiftRows(state []byte) {
	state[1], state[5], state[9], state[13] = state[5], state[9], state[13], state[1]
	state[2], state[6], state[10], state[14] = state[10], state[14], state[2], state[6]
	state[3], state[7], state[11], state[15] = state[15], state[3], state[7], state[11]
}
func InvShiftRows(state []byte) {
	state[1], state[5], state[9], state[13] = state[13], state[1], state[5], state[9]
	state[2], state[6], state[10], state[14] = state[10], state[14], state[2], state[6]
	state[3], state[7], state[11], state[15] = state[7], state[11], state[15], state[3]
}
func MixColumns(state []byte) {
	var Tmp, Tm, t byte
	for i := 0; i < 4; i++ {
		t = state[i*4]
		Tmp = state[i*4] ^ state[i*4+1] ^ state[i*4+2] ^ state[i*4+3]
		Tm = state[i*4] ^ state[i*4+1]
		Tm = gmul2(Tm)
		state[i*4] ^= Tm ^ Tmp
		Tm = state[i*4+1] ^ state[i*4+2]
		Tm = gmul2(Tm)
		state[i*4+1] ^= Tm ^ Tmp
		Tm = state[i*4+2] ^ state[i*4+3]
		Tm = gmul2(Tm)
		state[i*4+2] ^= Tm ^ Tmp
		Tm = state[i*4+3] ^ t
		Tm = gmul2(Tm)
		state[i*4+3] ^= Tm ^ Tmp

	}
}
func InvMixColumns(state []byte) {
	for i := 0; i < 4; i++ {
		a := state[i*4]
		b := state[i*4+1]
		c := state[i*4+2]
		d := state[i*4+3]

		state[i*4] = gmul(0x0e, a) ^ gmul(0x0b, b) ^ gmul(0x0d, c) ^ gmul(0x09, d)
		state[i*4+1] = gmul(0x09, a) ^ gmul(0x0e, b) ^ gmul(0x0b, c) ^ gmul(0x0d, d)
		state[i*4+2] = gmul(0x0d, a) ^ gmul(0x09, b) ^ gmul(0x0e, c) ^ gmul(0x0b, d)
		state[i*4+3] = gmul(0x0b, a) ^ gmul(0x0d, b) ^ gmul(0x09, c) ^ gmul(0x0e, d)
	}
}

func gmul2(x byte) byte {
	return ((x << 1) ^ (((x >> 7) & 1) * 0x1b))
}
func gmul(y, x byte) byte {
	return ((y & 1) * x) ^ ((y >> 1 & 1) * gmul2(x)) ^ ((y >> 2 & 1) * gmul2(gmul2(x))) ^ ((y >> 3 & 1) * gmul2(gmul2(gmul2(x)))) ^ ((y >> 4 & 1) * gmul2(gmul2(gmul2(gmul2(x)))))

}
func AES_init_ctx(ctx *AES_KEY, key []byte, In []byte) {
	for i := 0; i < 16; i++ {
		ctx.Iv[i] = In[i]
	}
	KeyExpansion(ctx.RoundKey[:], key)
}
func KeyExpansion(keyEx []byte, key []byte) {
	for i := 0; i < 4; i++ {
		keyEx[i*4] = key[i*4]
		keyEx[i*4+1] = key[i*4+1]
		keyEx[i*4+2] = key[i*4+2]
		keyEx[i*4+3] = key[i*4+3]
	}
	temp := [4]byte{}
	for i := 4; i < 4*11; i++ {
		k := (i - 1) * 4
		temp[0] = keyEx[k]
		temp[1] = keyEx[k+1]
		temp[2] = keyEx[k+2]
		temp[3] = keyEx[k+3]

		if i%4 == 0 {
			t := temp[0]
			temp[0] = stb[temp[1]]
			temp[1] = stb[temp[2]]
			temp[2] = stb[temp[3]]
			temp[3] = stb[t]
			temp[0] ^= RC[i/4]
		}
		j := i * 4
		k = (i - 4) * 4
		keyEx[j] = keyEx[k] ^ temp[0]
		keyEx[j+1] = keyEx[k+1] ^ temp[1]
		keyEx[j+2] = keyEx[k+2] ^ temp[2]
		keyEx[j+3] = keyEx[k+3] ^ temp[3]

	}
}
func XorWithIv(buf []byte, iv []byte) {

	for i := 0; i < 16; i++ {
		buf[i] ^= iv[i]
	}
}
func AddRoundKey(round int, state []byte, keyEx []byte) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			state[i*4+j] ^= keyEx[round*16+i*4+j]
		}
	}
}
func Cipher(state []byte, keyEx []byte) {
	AddRoundKey(0, state, keyEx)
	for i := 1; ; i++ {
		SubBytes(state)
		ShiftRows(state)

		if i == 10 {
			break
		}
		MixColumns(state)
		AddRoundKey(i, state, keyEx)
	}
	AddRoundKey(10, state, keyEx)
}
func InvCipher(state []byte, keyEx []byte) {
	AddRoundKey(10, state, keyEx)
	for i := 9; ; i-- {
		InvShiftRows(state)
		InvSubBytes(state)
		AddRoundKey(i, state, keyEx)
		if i == 0 {
			break
		}
		InvMixColumns(state)
	}
}
func AES_CBC_encrypt_buffer(ctx *AES_KEY, buf []byte, length int) {
	Iv := ctx.Iv[:]
	for i := 0; i < length; i += 16 {
		XorWithIv(buf[i:i+16], Iv[:])
		Cipher(buf[i:i+16], ctx.RoundKey[:])
		Iv = buf[i : i+16]
	}
	for i := 0; i < 16; i++ {
		ctx.Iv[i] = Iv[i]
	}
}
func AES_CBC_decrypt_buffer(ctx *AES_KEY, buf []byte, length int) {
	NextIv := make([]byte, 16)
	for i := 0; i < length; i += 16 {
		for j := 0; j < 16; j++ {
			NextIv[j] = buf[i+j]
		}
		InvCipher(buf[i:i+16], ctx.RoundKey[:])
		XorWithIv(buf[i:i+16], ctx.Iv[:])
		for j := 0; j < 16; j++ {
			ctx.Iv[j] = NextIv[j]
		}
	}
	// for i := 0; i < 16; i++ {
	// 	ctx.Iv[i] = Iv[i]
	// }
}
func AES_init_ctx_iv(ctx *AES_KEY, key []byte, iv []byte) {
	KeyExpansion(ctx.RoundKey[:], key)
	for i := 0; i < 16; i++ {
		ctx.Iv[i] = iv[i]
	}
}

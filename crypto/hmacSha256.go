package crypto

type Digest []uint32

const (
	hexBase = 16

	//HMAC
	blockSizeBytes = 64
	iPad           = 0x36
	oPad           = 0x5c

	//SHA256
	bits                = 8
	bytesInChunk        = 64
	bytesAdditional     = 5
	bytesInWord         = 4
	byteMaxShiftInWord  = 24
	wordsTotalInChunk   = 64
	wordsPrimaryInChunk = 16
	hashSizeBytes       = 32
)

var (
	//SHA256
	k = [64]uint32{
		0x428A2F98, 0x71374491, 0xB5C0FBCF, 0xE9B5DBA5, 0x3956C25B, 0x59F111F1, 0x923F82A4, 0xAB1C5ED5,
		0xD807AA98, 0x12835B01, 0x243185BE, 0x550C7DC3, 0x72BE5D74, 0x80DEB1FE, 0x9BDC06A7, 0xC19BF174,
		0xE49B69C1, 0xEFBE4786, 0x0FC19DC6, 0x240CA1CC, 0x2DE92C6F, 0x4A7484AA, 0x5CB0A9DC, 0x76F988DA,
		0x983E5152, 0xA831C66D, 0xB00327C8, 0xBF597FC7, 0xC6E00BF3, 0xD5A79147, 0x06CA6351, 0x14292967,
		0x27B70A85, 0x2E1B2138, 0x4D2C6DFC, 0x53380D13, 0x650A7354, 0x766A0ABB, 0x81C2C92E, 0x92722C85,
		0xA2BFE8A1, 0xA81A664B, 0xC24B8B70, 0xC76C51A3, 0xD192E819, 0xD6990624, 0xF40E3585, 0x106AA070,
		0x19A4C116, 0x1E376C08, 0x2748774C, 0x34B0BCB5, 0x391C0CB3, 0x4ED8AA4A, 0x5B9CCA4F, 0x682E6FF3,
		0x748F82EE, 0x78A5636F, 0x84C87814, 0x8CC70208, 0x90BEFFFA, 0xA4506CEB, 0xBEF9A3F7, 0xC67178F2}
)

func (h Digest) ToHexString() string {
	if h == nil {
		return ""
	}
	var b = make([]byte, hashSizeBytes*2, hashSizeBytes*2)
	for i := 0; i < len(h); i++ {
		v := h[i]
		for j := bytesInWord*2 - 1; j >= 0; j-- {
			b[j+i*bytesInWord*2] = toChar(v % hexBase)
			v /= hexBase
		}
	}
	return string(b)
}

func (h Digest) ToBytes() []byte {
	if h == nil {
		return nil
	}
	var b = make([]byte, hashSizeBytes, hashSizeBytes)
	for i := 0; i < len(b); i += 4 {
		b[i+0] = byte((h[i/4] & (^(uint32(1<<24) - 1))) >> 24)
		b[i+1] = byte((h[i/4] & (^(uint32(1<<16) - 1))) >> 16)
		b[i+2] = byte((h[i/4] & (^(uint32(1<<8) - 1))) >> 8)
		b[i+3] = byte((h[i/4] & (^(uint32(1<<0) - 1))) >> 0)
	}
	return b
}

func toChar(v uint32) byte {
	switch v {
	case 0:
		return '0'
	case 1:
		return '1'
	case 2:
		return '2'
	case 3:
		return '3'
	case 4:
		return '4'
	case 5:
		return '5'
	case 6:
		return '6'
	case 7:
		return '7'
	case 8:
		return '8'
	case 9:
		return '9'
	case 10:
		return 'a'
	case 11:
		return 'b'
	case 12:
		return 'c'
	case 13:
		return 'd'
	case 14:
		return 'e'
	case 15:
		return 'f'
	default:
		return 'g'
	}
}

func HMACSHA256(msg, key []byte) Digest {
	if len(key) == 0 {
		return nil
	}

	iPads := make([]byte, blockSizeBytes, blockSizeBytes)
	oPads := make([]byte, blockSizeBytes, blockSizeBytes)

	for i := 0; i < len(key); i++ {
		iPads[i] = key[i] ^ iPad
		oPads[i] = key[i] ^ oPad
	}
	for i := len(key); i < blockSizeBytes; i++ {
		iPads[i] = 0 ^ iPad
		oPads[i] = 0 ^ oPad
	}

	innerBlock := make([]byte, blockSizeBytes+len(msg), blockSizeBytes+len(msg))

	for i := 0; i < blockSizeBytes; i++ {
		innerBlock[i] = iPads[i]
	}
	for i, j := blockSizeBytes, 0; i < blockSizeBytes+len(msg); i, j = i+1, j+1 {
		innerBlock[i] = msg[j]
	}

	innerHash := SHA256(innerBlock).ToBytes()

	outerBlockSize := blockSizeBytes + hashSizeBytes
	outerBlock := make([]byte, outerBlockSize, outerBlockSize)
	for i := 0; i < blockSizeBytes; i++ {
		outerBlock[i] = oPads[i]
	}
	for i, j := blockSizeBytes, 0; i < outerBlockSize; i, j = i+1, j+1 {
		outerBlock[i] = innerHash[j]
	}

	return SHA256(outerBlock)
}

func SHA256(msg []byte) Digest {

	r := Digest{0x6A09E667, 0xBB67AE85, 0x3C6EF372, 0xA54FF53A, 0x510E527F, 0x9B05688C, 0x1F83D9AB, 0x5BE0CD19}

	oldLenBits := uint(len(msg) * bits)

	chunks := 0
	for bytesInChunk*chunks-len(msg)-bytesAdditional < 0 {
		chunks++
	}
	zeroBytes := bytesInChunk*chunks - len(msg) - bytesAdditional

	newLenBytes := len(msg) + bytesAdditional + zeroBytes

	w := make([]uint32, newLenBytes/bytesInWord, newLenBytes/bytesInWord)

	for i := 0; i < len(msg)/bytesInWord+1; i++ {
		for j := 0; j < bytesInWord && j+i*bytesInWord < len(msg); j++ {
			w[i] += uint32(msg[j+i*bytesInWord]) << (byteMaxShiftInWord - j*bits)
		}
	}

	w[len(msg)/bytesInWord] |= (1 << (bits - 1)) << (bits * (bytesInWord - 1 - len(msg)%bytesInWord))
	w[len(w)-1] = uint32(oldLenBits)

	for i := 0; i < len(w)/wordsPrimaryInChunk; i++ {
		wAll := make([]uint32, wordsTotalInChunk, wordsTotalInChunk)
		for j := 0; j < wordsPrimaryInChunk; j++ {
			wAll[j] = w[j+i*wordsPrimaryInChunk]
		}
		for j := wordsPrimaryInChunk; j < wordsTotalInChunk; j++ {
			s0 := RCRU32(wAll[j-15], 7) ^ RCRU32(wAll[j-15], 18) ^ (wAll[j-15] >> 3)
			s1 := RCRU32(wAll[j-2], 17) ^ RCRU32(wAll[j-2], 19) ^ (wAll[j-2] >> 10)
			wAll[j] = wAll[j-16] + s0 + wAll[j-7] + s1
		}
		a, b, c, d, e, f, g, h := r[0], r[1], r[2], r[3], r[4], r[5], r[6], r[7]
		for j := 0; j < wordsTotalInChunk; j++ {
			sum0 := RCRU32(a, 2) ^ RCRU32(a, 13) ^ RCRU32(a, 22)
			Ma := (a & b) ^ (a & c) ^ (b & c)
			t2 := sum0 + Ma
			sum1 := RCRU32(e, 6) ^ RCRU32(e, 11) ^ RCRU32(e, 25)
			Ch := (e & f) ^ ((^e) & g)
			t1 := h + sum1 + Ch + k[j] + wAll[j]
			h, g, f, e, d, c, b, a = g, f, e, d+t1, c, b, a, t1+t2
		}
		r[0], r[1], r[2], r[3], r[4], r[5], r[6], r[7] = r[0]+a, r[1]+b, r[2]+c, r[3]+d, r[4]+e, r[5]+f, r[6]+g, r[7]+h
	}

	return r
}

func RCRU32(v uint32, i byte) uint32 {
	return (v >> i) + ((v & ((1 << i) - 1)) << (32 - i))
}

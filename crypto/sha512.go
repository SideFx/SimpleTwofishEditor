//----------------------------------------------------------------------------------------------------------------------
// (w) 2022-2024 by Jan Buchholz
// Sha512 Go port based on Wikipedia's pseudo code
// https://en.wikipedia.org/wiki/SHA-2
//----------------------------------------------------------------------------------------------------------------------

package crypto

const Sha512Shabytes int = 64
const uint64Bits int = 64
const shaRounds uint64 = 80
const shaBlockSize uint64 = 128
const shaReservedBytes uint64 = 16

type ShaResult [Sha512Shabytes]byte

type Sha512 struct {
	d0     [shaRounds]uint64
	h0     [8]uint64
	t0     [8]uint64
	k0     [shaRounds]uint64
	size   uint64
	result ShaResult
}

func NewSha512() Sha512 {
	sha := Sha512{}
	sha.k0 = [...]uint64{
		0x428a2f98d728ae22, 0x7137449123ef65cd, 0xb5c0fbcfec4d3b2f, 0xe9b5dba58189dbbc,
		0x3956c25bf348b538, 0x59f111f1b605d019, 0x923f82a4af194f9b, 0xab1c5ed5da6d8118,
		0xd807aa98a3030242, 0x12835b0145706fbe, 0x243185be4ee4b28c, 0x550c7dc3d5ffb4e2,
		0x72be5d74f27b896f, 0x80deb1fe3b1696b1, 0x9bdc06a725c71235, 0xc19bf174cf692694,
		0xe49b69c19ef14ad2, 0xefbe4786384f25e3, 0x0fc19dc68b8cd5b5, 0x240ca1cc77ac9c65,
		0x2de92c6f592b0275, 0x4a7484aa6ea6e483, 0x5cb0a9dcbd41fbd4, 0x76f988da831153b5,
		0x983e5152ee66dfab, 0xa831c66d2db43210, 0xb00327c898fb213f, 0xbf597fc7beef0ee4,
		0xc6e00bf33da88fc2, 0xd5a79147930aa725, 0x06ca6351e003826f, 0x142929670a0e6e70,
		0x27b70a8546d22ffc, 0x2e1b21385c26c926, 0x4d2c6dfc5ac42aed, 0x53380d139d95b3df,
		0x650a73548baf63de, 0x766a0abb3c77b2a8, 0x81c2c92e47edaee6, 0x92722c851482353b,
		0xa2bfe8a14cf10364, 0xa81a664bbc423001, 0xc24b8b70d0f89791, 0xc76c51a30654be30,
		0xd192e819d6ef5218, 0xd69906245565a910, 0xf40e35855771202a, 0x106aa07032bbd1b8,
		0x19a4c116b8d2d0c8, 0x1e376c085141ab53, 0x2748774cdf8eeb99, 0x34b0bcb5e19b48a8,
		0x391c0cb3c5c95a63, 0x4ed8aa4ae3418acb, 0x5b9cca4f7763e373, 0x682e6ff3d6b2b8a3,
		0x748f82ee5defb2fc, 0x78a5636f43172f60, 0x84c87814a1f0ab72, 0x8cc702081a6439ec,
		0x90befffa23631e28, 0xa4506cebde82bde9, 0xbef9a3f7b2c67915, 0xc67178f2e372532b,
		0xca273eceea26619c, 0xd186b8c721c0c207, 0xeada7dd6cde0eb1e, 0xf57d4f7fee6ed178,
		0x06f067aa72176fba, 0x0a637dc5a2c898a6, 0x113f9804bef90dae, 0x1b710b35131c471b,
		0x28db77f523047d84, 0x32caab7b40c72493, 0x3c9ebe0a15c9bebc, 0x431d67c49c100d4c,
		0x4cc5d4becb3e42b6, 0x597f299cfc657e2a, 0x5fcb6fab3ad6faec, 0x6c44198c4a475817}
	return sha
}

func (sha Sha512) Compute(p []byte) ShaResult {
	var ha, hb, hc, hd, he, hf, hg, hh, sl, sh, sc, sm, t0, t1, z, n, i, j uint64
	sha.h0 = [...]uint64{
		0x6a09e667f3bcc908, 0xbb67ae8584caa73b, 0x3c6ef372fe94f82b, 0xa54ff53a5f1d36f1,
		0x510e527fade682d1, 0x9b05688c2b3e6c1f, 0x1f83d9abfb41bd6b, 0x5be0cd19137e2179}
	sha.size = uint64(len(p))
	var sd = sha.size/shaBlockSize + 1
	if (sha.size % shaBlockSize) >= (shaBlockSize - shaReservedBytes) {
		sd++
	}
	// get number of bits
	sl, sh = sha.size<<3, sha.size>>61
	mc := make([]byte, sha.size)
	copy(mc, p)
	mc = append(mc, 0x80)
	for i = sha.size + 1; i < sd*shaBlockSize-shaReservedBytes; i++ {
		mc = append(mc, 0x00)
	}
	mc = append(mc, byte(sh>>56))
	mc = append(mc, byte(sh>>48))
	mc = append(mc, byte(sh>>40))
	mc = append(mc, byte(sh>>32))
	mc = append(mc, byte(sh>>24))
	mc = append(mc, byte(sh>>16))
	mc = append(mc, byte(sh>>8))
	mc = append(mc, byte(sh))
	mc = append(mc, byte(sl>>56))
	mc = append(mc, byte(sl>>48))
	mc = append(mc, byte(sl>>40))
	mc = append(mc, byte(sl>>32))
	mc = append(mc, byte(sl>>24))
	mc = append(mc, byte(sl>>16))
	mc = append(mc, byte(sl>>8))
	mc = append(mc, byte(sl))
	z = 0
	for n = 0; n < sd; n++ {
		ha = sha.h0[0]
		hb = sha.h0[1]
		hc = sha.h0[2]
		hd = sha.h0[3]
		he = sha.h0[4]
		hf = sha.h0[5]
		hg = sha.h0[6]
		hh = sha.h0[7]
		for i = 0; i < 16; i++ {
			for j = 0; j < 8; j++ {
				sha.t0[j] = uint64(mc[j+z])
			}
			sha.d0[i] = (sha.t0[0] << 56) | (sha.t0[1] << 48) | (sha.t0[2] << 40) | (sha.t0[3] << 32) |
				(sha.t0[4] << 24) | (sha.t0[5] << 16) | (sha.t0[6] << 8) | sha.t0[7]
			z += 8
		}
		for i = 16; i < shaRounds; i++ {
			sl = ror64(sha.d0[i-15], 1) ^ ror64(sha.d0[i-15], 8) ^ (sha.d0[i-15] >> 7)
			sh = ror64(sha.d0[i-2], 19) ^ ror64(sha.d0[i-2], 61) ^ (sha.d0[i-2] >> 6)
			sha.d0[i] = sha.d0[i-16] + sl + sha.d0[i-7] + sh
		}
		for i = 0; i < shaRounds; i++ {
			sh = ror64(he, 14) ^ ror64(he, 18) ^ ror64(he, 41)
			sc = (he & hf) ^ ((^he) & hg)
			t0 = hh + sh + sc + sha.k0[i] + sha.d0[i]
			sl = ror64(ha, 28) ^ ror64(ha, 34) ^ ror64(ha, 39)
			sm = (ha & hb) ^ (ha & hc) ^ (hb & hc)
			t1 = sl + sm
			hh = hg
			hg = hf
			hf = he
			he = hd + t0
			hd = hc
			hc = hb
			hb = ha
			ha = t0 + t1
		}
		sha.h0[0] += ha
		sha.h0[1] += hb
		sha.h0[2] += hc
		sha.h0[3] += hd
		sha.h0[4] += he
		sha.h0[5] += hf
		sha.h0[6] += hg
		sha.h0[7] += hh
	}
	i = 0
	for j = 0; j < 8; j++ {
		sha.result[i] = byte(sha.h0[j] >> 56)
		sha.result[i+1] = byte(sha.h0[j] >> 48)
		sha.result[i+2] = byte(sha.h0[j] >> 40)
		sha.result[i+3] = byte(sha.h0[j] >> 32)
		sha.result[i+4] = byte(sha.h0[j] >> 24)
		sha.result[i+5] = byte(sha.h0[j] >> 16)
		sha.result[i+6] = byte(sha.h0[j] >> 8)
		sha.result[i+7] = byte(sha.h0[j])
		i += 8
	}
	return sha.result
}

func ror64(x uint64, n int) uint64 {
	return (x >> n) | (x << (uint64Bits - n))
}

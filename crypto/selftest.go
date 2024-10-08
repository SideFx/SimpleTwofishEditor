//----------------------------------------------------------------------------------------------------------------------
// (w) 2023-2024 by Jan Buchholz
// Self tests for Sha512, Blowfish and Twofish
//----------------------------------------------------------------------------------------------------------------------

package crypto

import "slices"

func SelfTest() bool {
	return checkTwofish() && checkSha512()
}

func checkTwofish() bool {
	var ct49 = TfBlock{
		0x37, 0xfe, 0x26, 0xff, 0x1c, 0xf6, 0x61, 0x75, 0xf5, 0xdd, 0xf4, 0xc3, 0x3b, 0x97, 0xa2, 0x05}
	const rounds uint32 = 49
	var tfkey TfKey
	var pt, ptx, ct TfBlock
	var i, j uint32
	var tfish Twofish
	for i = 0; i < TwofishKeysize; i++ {
		tfkey[i] = 0x00
	}
	for i = 0; i < TwofishBlocksize; i++ {
		pt[i] = 0x00
	}
	for i = 0; i < rounds; i++ {
		tfish = NewTwofish(tfkey)
		copy(ptx[:], pt[:])
		tfish.EncryptBlock(&pt)
		copy(ct[:], pt[:])
		tfish.DecryptBlock(&pt)
		if !slices.Equal(pt[:], ptx[:]) {
			return false
		}
		for j = TwofishKeysize / 2; j < TwofishKeysize; j++ {
			tfkey[j] = tfkey[j-TwofishKeysize/2]
		}
		for j = 0; j < TwofishBlocksize; j++ {
			tfkey[j] = ptx[j]
		}
		copy(pt[:], ct[:])
	}
	return slices.Equal(ct[:], ct49[:])
}

func checkSha512() bool {
	var check = ShaResult{
		0x07, 0xe5, 0x47, 0xd9, 0x58, 0x6f, 0x6a, 0x73, 0xf7, 0x3f, 0xba, 0xc0, 0x43, 0x5e, 0xd7, 0x69,
		0x51, 0x21, 0x8f, 0xb7, 0xd0, 0xc8, 0xd7, 0x88, 0xa3, 0x09, 0xd7, 0x85, 0x43, 0x6b, 0xbb, 0x64,
		0x2e, 0x93, 0xa2, 0x52, 0xa9, 0x54, 0xf2, 0x39, 0x12, 0x54, 0x7d, 0x1e, 0x8a, 0x3b, 0x5e, 0xd6,
		0xe1, 0xbf, 0xd7, 0x09, 0x78, 0x21, 0x23, 0x3f, 0xa0, 0x53, 0x8f, 0x3d, 0xb8, 0x54, 0xfe, 0xe6}
	inp := []byte("The quick brown fox jumps over the lazy dog")
	sha := NewSha512()
	res := sha.Compute(inp)
	return slices.Equal(check[:], res[:])
}

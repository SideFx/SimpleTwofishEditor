//----------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz
// Create final encryption/decryption key from password entered
//----------------------------------------------------------------------------------------------------------------------

package crypto

const shaKeyRounds = 1234
const bits = 8

var vault [TwofishKeysize]byte
var valid bool

func init() {
	Invalidate()
}

func Push(p []byte) {
	var i int
	var l uint32
	var tmp TfKey
	sha := NewSha512()
	buffer := sha.Compute(p)
	for i = 0; i < shaKeyRounds; i++ {
		buffer = sha.Compute(buffer[:])
	}
	for l = 0; l < TwofishKeysize; l++ {
		tmp[l] = buffer[l] + buffer[l+TwofishKeysize]
	}
	vault = encode(tmp)
	Validate()
}

func Pop() TfKey {
	return decode(vault)
}

func Invalidate() {
	vault = [TwofishKeysize]byte{}
	valid = false
}

func Validate() {
	valid = true
}

func IsValid() bool {
	return valid
}

func encode(b TfKey) TfKey {
	var c TfKey
	copy(c[:], b[:])
	j := 0
	for i, e := range c {
		i++
		if i > bits-1 {
			i = 1
		}
		c[j] = ror(e, i)
		j++
	}
	return c
}

func decode(b TfKey) TfKey {
	var c TfKey
	copy(c[:], b[:])
	j := 0
	for i, d := range c {
		i++
		if i > bits-1 {
			i = 1
		}
		c[j] = rol(d, i)
		j++
	}
	return c
}

func ror(x byte, n int) byte {
	return (x >> n) | (x << (bits - n))
}

func rol(x byte, n int) byte {
	return (x << n) | (x >> (bits - n))
}

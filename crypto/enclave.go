//----------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz
// Create final encryption/decryption key from password entered
//----------------------------------------------------------------------------------------------------------------------

package crypto

var vault [TwofishKeysize]byte
var valid bool

func init() {
	Invalidate()
}

func Push(p []byte) {
	var tmp TfKey
	var lp, lv = len(p), len(tmp)
	var i, j int
	if lp == lv { // equal size
		copy(vault[:], p)
	} else {
		if lp > lv { // password longer than vault
			j = 0
			for i = 0; i < lv; i++ {
				tmp[i] = p[i]
			}
			for i = lv; i < lp; i++ {
				tmp[j] = tmp[j] + p[i]
				j++
			}
		} else { // vault longer than password
			j = 0
			for i = 0; i < lv; i++ {
				if j >= lp {
					j = 0
				}
				tmp[i] = p[j]
				j++
			}
		}
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

const bits = 8

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

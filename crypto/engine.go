//----------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz
// Encrypt/decrypt text
//----------------------------------------------------------------------------------------------------------------------

package crypto

import (
	"SimpleTwofishEditor/assets"
	"crypto/rand"
	"slices"
)

const tokenSize = 16

var dataPrefix = []byte("!SiMpLe!TwOfIsH!EdItOr!")

func EncryptPayload(payload []byte) ([]byte, error) {
	var outp = make([]byte, len(dataPrefix))
	copy(outp, dataPrefix)
	if len(payload) > 0 {
		data := make([]byte, len(payload))
		copy(data, payload)
		token := make([]byte, tokenSize)
		if _, err := rand.Read(token); err != nil {
			return nil, err
		}
		stage0 := make([]byte, tokenSize)
		copy(stage0, token)
		stage0 = append(stage0, data...)
		sha := NewSha512()
		shaResult := sha.Compute(stage0)
		outp = append(outp, shaResult[:]...)
		tf := NewTwofishWithEnclave()
		stage1 := tf.CbcEncrypt(stage0)
		outp = append(outp, stage1...)
	}
	return outp, nil
}

func DecryptPayload(payload []byte) (string, string) {
	if len(payload) > 0 {
		data := make([]byte, len(payload))
		copy(data, payload)
		if len(data) < len(dataPrefix) {
			return "", assets.ErrNoMatch
		}
		for i := 0; i < len(dataPrefix); i++ {
			if dataPrefix[i] != data[i] {
				return "", assets.ErrNoMatch
			}
		}
		if len(data) == len(dataPrefix) {
			return "", "" //empty Zydeco file
		}
		if len(data) < tokenSize+len(dataPrefix)+Sha512Shabytes+1 {
			return "", assets.ErrCorrupted
		}
		data = data[len(dataPrefix):]
		shaCheck := data[:Sha512Shabytes]
		data = data[Sha512Shabytes:]
		tf := NewTwofishWithEnclave()
		tmp := tf.CbcDecrypt(data)
		sha := NewSha512()
		shaResult := sha.Compute(tmp)
		tmp = tmp[tokenSize:]
		r := slices.Equal(shaCheck[:], shaResult[:])
		if !r {
			return "", assets.ErrUnableToDecrypt
		}
		return string(tmp), ""
	}
	return "", assets.ErrEmptyFile
}

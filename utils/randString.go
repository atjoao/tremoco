//  https://go.dev/play/p/1GwSRsKIsd

package utils

import (
	"crypto/rand"
	"math/big"
	"unicode"
)

func RandString(n int) string {
	max := big.NewInt(130)
	bs := make([]byte, n)

	for i := range bs {
		g, _ := rand.Int(rand.Reader, max)
		r := rune(g.Int64())
		for !unicode.IsNumber(r) && !unicode.IsLetter(r) {
			g, _ = rand.Int(rand.Reader, max)
			r = rune(g.Int64())
		}
		bs[i] = byte(g.Int64())
	}
	return string(bs)
}

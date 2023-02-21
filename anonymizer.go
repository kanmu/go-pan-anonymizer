package anonymizer

import (
	"golang.org/x/text/transform"
)

type anonymizer struct {
	transform.NopResetter
	testPan func([]byte) bool
	mask    byte
}

func DefaultAnonymizer() *anonymizer {
	return &anonymizer{
		NopResetter: transform.NopResetter{},
		testPan:     TestLuhn,
		mask:        '*',
	}
}

func NewAnonymizer(mask byte, test func([]byte) bool) *anonymizer {
	return &anonymizer{
		NopResetter: transform.NopResetter{},
		testPan:     test,
		mask:        mask,
	}
}

func (a *anonymizer) Transform(dst, src []byte, atEOF bool) (int, int, error) {
	n := 0
	nDst := 0
	nSrc := 0
	for i := 0; i < len(src); i++ {
		b := src[i]
		if b >= '0' && b <= '9' {
			if n < 16 {
				n += 1
			}
		} else {
			n = 0
		}
		if n == 16 {
			mf := i - 15
			if a.testPan(src[mf : i+1]) {
				if nSrc <= mf {
					nDst += copy(dst[nDst:], src[nSrc:mf])
				}
				for ; nDst <= i; nDst++ {
					dst[nDst] = a.mask
				}
				nSrc = nDst
			}
		}
	}
	if !atEOF && n > 0 && n < 16 {
		nDst += copy(dst[nDst:], src[nSrc:len(src)-n])
		nSrc = nDst
		return nDst, nSrc, transform.ErrShortSrc
	}
	nDst += copy(dst[nDst:], src[nSrc:])
	nSrc = nDst
	return nDst, nSrc, nil
}

func TestLuhn(digits []byte) bool {
	sum := 0
	for i := 0; i < len(digits); i++ {
		d := int(digits[i] - '0')
		if i%2 == 0 {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
	}
	return sum%10 == 0
}

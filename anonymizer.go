package anonymizer

import (
	"golang.org/x/text/transform"
)

type Anonymizer struct {
	transform.NopResetter
	testPan func([]byte) bool
}

func DefaultAnonymizer() *Anonymizer {
	return &Anonymizer{
		testPan: TestLuhn,
	}
}

func NewAnonymizer(test func([]byte) bool) *Anonymizer {
	return &Anonymizer{
		testPan: test,
	}
}

func (a *Anonymizer) Transform(dst, src []byte, atEOF bool) (int, int, error) {
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
					nDst += copy(dst[nDst:], "****************")
				} else {
					for ; nDst < i; nDst++ {
						dst[nDst] = '*'
					}
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

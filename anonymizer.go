package anonymizer

import (
	"golang.org/x/text/transform"
)

type Anonymizer struct {
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
	i := 0
	n := 0
	nDst := 0
	nSrc := 0
	for ; i < len(src); i++ {
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
	n := len(digits)
	sum := 0
	alt := true
	for i := 0; i < n; i++ {
		d := int(digits[i] - 0x30)
		if alt {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
		alt = !alt
	}
	return sum%10 == 0
}

func (a *Anonymizer) Reset() {
}

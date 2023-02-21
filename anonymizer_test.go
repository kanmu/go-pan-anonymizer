package anonymizer

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/transform"
)

func TestTransform(t *testing.T) {
	test := func(src string, expected string) {
		r := strings.NewReader(src)
		a := DefaultAnonymizer()
		tf := transform.NewReader(r, a)
		out, err := io.ReadAll(tf)
		assert.NoError(t, err)
		assert.Equal(t, expected, string(out))
	}

	t.Run("in alpha", func(t *testing.T) {
		src := "abcdef4242424242424242ghijk"
		test(src, "abcdef****************ghijk")
	})
	t.Run("in digits", func(t *testing.T) {
		src := "9994242424242424242999"
		test(src, "999****************999")
	})
	t.Run("end", func(t *testing.T) {
		src := "abcdef4242424242424242"
		test(src, "abcdef****************")
	})
	t.Run("sep by digits", func(t *testing.T) {
		src := "994242424242424242999424242424242424299"
		test(src, "99****************999****************99")
	})
	t.Run("sep by alpha", func(t *testing.T) {
		src := "AA4242424242424242TTT4242424242424242AAA"
		test(src, "AA****************TTT****************AAA")
	})
	t.Run("connectted", func(t *testing.T) {
		// "42424242424242424242"[:16] and [2:18] is passed testLuhn
		src := "994242424242424242424299"
		test(src, "99********************99")
	})
	t.Run("large", func(t *testing.T) {
		src := "@@4242424242424242@@" + strings.Repeat("@", 4098)
		test(src, strings.ReplaceAll(src, "4242424242424242", "****************"))
	})
	t.Run("ErrShortSrc", func(t *testing.T) {
		src := strings.Repeat("@", 4094) + "4242424242424242@@"
		test(src, strings.ReplaceAll(src, "4242424242424242", "****************"))
	})
	t.Run("ErrShortSrc 2", func(t *testing.T) {
		src := strings.Repeat("@", 4096-16) + "424242424242424277"
		test(src, strings.ReplaceAll(src, "4242424242424242", "****************"))
	})
	t.Run("small", func(t *testing.T) {
		src := "123"
		test(src, "123")
	})
}

func TestTransform_WithCustomTester(t *testing.T) {
	t.Run("hit", func(t *testing.T) {
		a := NewAnonymizer('*', func(bs []byte) bool {
			return strings.HasPrefix(string(bs), "0123")
		})
		src := "!0123456789012345!"
		r := strings.NewReader(src)
		tf := transform.NewReader(r, a)
		out, err := io.ReadAll(tf)
		assert.NoError(t, err)
		assert.Equal(t, "!****************!", string(out))
	})
	t.Run("miss", func(t *testing.T) {
		a := NewAnonymizer('*', func(bs []byte) bool {
			return strings.HasPrefix(string(bs), "0123")
		})
		src := "!9999999999999999!"
		r := strings.NewReader(src)
		tf := transform.NewReader(r, a)
		out, err := io.ReadAll(tf)
		assert.NoError(t, err)
		assert.Equal(t, "!9999999999999999!", string(out))
	})
	t.Run("mask char", func(t *testing.T) {
		a := NewAnonymizer('?', func(bs []byte) bool {
			return strings.HasPrefix(string(bs), "0123")
		})
		src := "!0123456789012345!"
		r := strings.NewReader(src)
		tf := transform.NewReader(r, a)
		out, err := io.ReadAll(tf)
		assert.NoError(t, err)
		assert.Equal(t, "!????????????????!", string(out))
	})
}

func TestTestLuhn(t *testing.T) {
	assert.True(t, TestLuhn([]byte("4242424242424242")))
	assert.False(t, TestLuhn([]byte("4242424242424241")))
}

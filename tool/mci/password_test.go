package mci

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePassword(t *testing.T) {
	const l = 16
	s := GeneratePassword(l)
	assert.Len(t, s, l)

	s = GeneratePassword(l, Upper, Lower, Digits)
	assert.Len(t, s, l)

	s = GeneratePassword(l, Upper, Lower, Digits)
	assert.Len(t, s, l)

	s = GeneratePasswordBySet(l, UpperLetters, DigitsLetters, LowerLetters, "-#")
	assert.Len(t, s, l)
}

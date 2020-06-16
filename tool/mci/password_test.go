package mci_test

import (
	"testing"

	"github.com/bingoohuang/tool/mci"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePassword(t *testing.T) {
	const l = 16
	s := mci.GeneratePassword(l)
	assert.Len(t, s, l)

	s = mci.GeneratePassword(l, mci.Upper, mci.Lower, mci.Digits)
	assert.Len(t, s, l)

	s = mci.GeneratePassword(l, mci.Upper, mci.Lower, mci.Digits)
	assert.Len(t, s, l)

	s = mci.GeneratePasswordBySet(l, mci.UpperLetters, mci.DigitsLetters, mci.LowerLetters, "-#")
	assert.Len(t, s, l)
}

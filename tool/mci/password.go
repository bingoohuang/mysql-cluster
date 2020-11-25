package mci

import (
	"crypto/rand"
	"math/big"
)

const (
	// LowerLetters is the list of lowercase letters.
	LowerLetters = "abcdefghijklmnopqrstuvwxyz"

	// UpperLetters is the list of uppercase letters.
	UpperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// DigitsLetters is the list of permitted digits.
	DigitsLetters = "0123456789"

	// SymbolsLetters is the list of symbols.
	SymbolsLetters = "~!@#$%^&*()_+`-={}|[]\\:\"<>?,./"
)

// PasswordType means password type.
type PasswordType int

const (
	// Lower lower letters.
	Lower PasswordType = iota
	// Upper upper letters.
	Upper
	// Digits digits.
	Digits
	// Symbols symbols.
	Symbols
)

// GeneratePasswordBySet generates a password with length n by types.
func GeneratePasswordBySet(n int, sets ...string) string {
	p, err := GeneratePasswordBySetE(n, sets...)
	if err != nil {
		panic(err)
	}

	return p
}

// GeneratePasswordBySetE generates a password with length n by types by specified sets.
func GeneratePasswordBySetE(n int, sets ...string) (string, error) {
	s := ""

	if len(sets) == 0 {
		sets = []string{LowerLetters, UpperLetters, DigitsLetters}
	}

	l := len(sets)

	for i := 0; i < n && i < l; i++ {
		r, err := randomElement(sets[i])
		if err != nil {
			return s, err
		}

		s += r
	}

	for i := l; i < n; i++ {
		ri, err := rand.Int(rand.Reader, big.NewInt(int64(l)))
		if err != nil {
			return "", err
		}

		r, err := randomElement(sets[ri.Int64()])
		if err != nil {
			return s, err
		}

		s += r
	}

	return s, nil
}

// GeneratePassword generates a password with length n by types.
func GeneratePassword(n int, types ...PasswordType) string {
	p, err := GeneratePasswordE(n, types...)
	if err != nil {
		panic(err)
	}

	return p
}

// GeneratePasswordE generates a password with length n by types.
func GeneratePasswordE(n int, types ...PasswordType) (string, error) {
	sets := make([]string, len(types))
	for i, t := range types {
		sets[i] = typeLetter(t)
	}

	return GeneratePasswordBySetE(n, sets...)
}

func typeLetter(pt PasswordType) string {
	switch pt {
	case Lower:
		return LowerLetters
	case Upper:
		return UpperLetters
	case Digits:
		return DigitsLetters
	case Symbols:
		return SymbolsLetters
	default:
		return LowerLetters
	}
}

// randomElement extracts a random element from the given string.
func randomElement(s string) (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(s))))
	if err != nil {
		return "", err
	}

	return string(s[n.Int64()]), nil
}

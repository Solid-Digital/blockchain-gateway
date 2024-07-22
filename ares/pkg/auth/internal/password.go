package internal

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/unchainio/pkg/errors"
)

const (
	// LowerLetters is the list of lowercase letters.
	LowerLetters = "abcdefghijklmnopqrstuvwxyz"

	// UpperLetters is the list of uppercase letters.
	UpperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Digits is the list of permitted digits.
	Digits = "0123456789"

	// Symbols is the list of symbols.
	Symbols = "~!@#$%^&*()_+`-={}|[]\\:\"<>?,./"

	// Minimum length
	MinLength = 8
)

var (
	// ErrExceedsTotalLength is the error returned with the number of digits and
	// symbols is greater than the total length.
	ErrExceedsTotalLength = errors.New("number of digits and symbols must be less than total length")

	// ErrLettersExceedsAvailable is the error returned with the number of letters
	// exceeds the number of available letters and repeats are not allowed.
	ErrLettersExceedsAvailable = errors.New("number of letters exceeds available letters and repeats are not allowed")

	// ErrDigitsExceedsAvailable is the error returned with the number of digits
	// exceeds the number of available digits and repeats are not allowed.
	ErrDigitsExceedsAvailable = errors.New("number of digits exceeds available digits and repeats are not allowed")

	// ErrSymbolsExceedsAvailable is the error returned with the number of symbols
	// exceeds the number of available symbols and repeats are not allowed.
	ErrSymbolsExceedsAvailable = errors.New("number of symbols exceeds available symbols and repeats are not allowed")

	// ErrInvalidLength is the error returned with the password length below minimum
	ErrInvalidLength = errors.New(fmt.Sprintf("password must be at least %d characters", MinLength))

	// ErrContainsNoDigits is the error returned with the password contains no digits
	ErrContainsNoDigits = errors.New("password must contain at least one digit")

	// ErrContainsNoUpperLetters is the error returned with the password contains no upper letters
	ErrContainsNoUpperLetters = errors.New("password must contain at least one upper letter")

	// ErrContainsNoLowerLetters is the error returned with the password contains no lower letters
	ErrContainsNoLowerLetters = errors.New("password must contain at least one lower letter")

	// ErrContainsNoSymbols is the error returned with the password contains no symbols
	ErrContainsNoSymbols = errors.New("password must contain at least one symbol")
)

// Generates a password with the given requirements. length is the
// total number of characters in the password. numDigits is the number of digits
// to include in the result. numSymbols is the number of symbols to include in
// the result. noUpper excludes uppercase letters from the results. allowRepeat
// allows characters to repeat.
func GeneratePassword(length, numDigits, numSymbols int, noUpper, allowRepeat bool) (string, error) {
	letters := LowerLetters
	if !noUpper {
		letters += UpperLetters
	}

	chars := length - numDigits - numSymbols
	if chars < 0 {
		return "", ErrExceedsTotalLength
	}

	if !allowRepeat && chars > len(letters) {
		return "", ErrLettersExceedsAvailable
	}

	if !allowRepeat && numDigits > len(Digits) {
		return "", ErrDigitsExceedsAvailable
	}

	if !allowRepeat && numSymbols > len(Symbols) {
		return "", ErrSymbolsExceedsAvailable
	}

	var result string

	// Characters
	for i := 0; i < chars; i++ {
		ch, err := randomElement(letters)
		if err != nil {
			return "", err
		}

		if !allowRepeat && strings.Contains(result, ch) {
			i--
			continue
		}

		result, err = randomInsert(result, ch)
		if err != nil {
			return "", err
		}
	}

	// Digits
	for i := 0; i < numDigits; i++ {
		d, err := randomElement(Digits)
		if err != nil {
			return "", err
		}

		if !allowRepeat && strings.Contains(result, d) {
			i--
			continue
		}

		result, err = randomInsert(result, d)
		if err != nil {
			return "", err
		}
	}

	// Symbols
	for i := 0; i < numSymbols; i++ {
		sym, err := randomElement(Symbols)
		if err != nil {
			return "", err
		}

		if !allowRepeat && strings.Contains(result, sym) {
			i--
			continue
		}

		result, err = randomInsert(result, sym)
		if err != nil {
			return "", err
		}
	}

	return result, nil
}

// randomInsert randomly inserts the given value into the given string.
func randomInsert(s, val string) (string, error) {
	if s == "" {
		return val, nil
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(s))))
	if err != nil {
		return "", err
	}
	i := n.Int64()
	return s[0:i] + val + s[i:len(s)], nil
}

// randomElement extracts a random element from the given string.
func randomElement(s string) (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(s))))
	if err != nil {
		return "", err
	}
	return string(s[n.Int64()]), nil
}

func validatePassword(password string) (bool, error) {

	if len(password) < MinLength {
		return false, ErrInvalidLength
	}

	if !strings.Contains(password, Digits) {
		return false, ErrContainsNoDigits
	}

	if !strings.Contains(password, UpperLetters) {
		return false, ErrContainsNoUpperLetters
	}

	if !strings.Contains(password, LowerLetters) {
		return false, ErrContainsNoLowerLetters
	}

	if !strings.Contains(password, Symbols) {
		return false, ErrContainsNoSymbols
	}

	return true, nil
}

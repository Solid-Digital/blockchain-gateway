package internal

import (
	"testing"
)

func TestPasswords(t *testing.T) {
	tr := initTestRunner(t)
	t.Parallel()

	tr.TestGeneratePassword(t)
}

type TestRunner struct {
	t *testing.T
}

func initTestRunner(t *testing.T) *TestRunner {
	return &TestRunner{}
}

func (tr *TestRunner) hasDuplicates(s string) bool {
	found := make(map[rune]struct{}, len(s))
	for _, ch := range s {
		if _, ok := found[ch]; ok {
			return true
		}
		found[ch] = struct{}{}
	}
	return false
}

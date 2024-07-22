/*
Copied and modified from https://github.com/stretchr/testify to make the assert.Error and require.Error work with custom errors
*/
package xrequire

import (
	"bitbucket.org/unchain/ares/pkg/testhelper/xassert"
)

// NoError asserts that a function returned no error (i.e. `nil`).
//
//   actualObj, err := SomeFunction()
//   if assert.NoError(t, err) {
// 	   assert.Equal(t, expectedObj, actualObj)
//   }
func NoError(t TestingT, err error, msgAndArgs ...interface{}) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if xassert.NoError(t, err, msgAndArgs...) {
		return
	}
	t.FailNow()
}

// Error asserts that a function returned an error (i.e. not `nil`).
//
//   actualObj, err := SomeFunction()
//   if assert.Error(t, err) {
// 	   assert.Equal(t, expectedError, err)
//   }
func Error(t TestingT, err error, msgAndArgs ...interface{}) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if xassert.Error(t, err, msgAndArgs...) {
		return
	}
	t.FailNow()
}

// Errorf asserts that a function returned an error (i.e. not `nil`).
//
//   actualObj, err := SomeFunction()
//   if assert.Errorf(t, err, "error message %s", "formatted") {
// 	   assert.Equal(t, expectedErrorf, err)
//   }
func Errorf(t TestingT, err error, msg string, args ...interface{}) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if xassert.Errorf(t, err, msg, args...) {
		return
	}
	t.FailNow()
}

// TestingT is an interface wrapper around *testing.T
type TestingT interface {
	Errorf(format string, args ...interface{})
	FailNow()
}

type tHelper interface {
	Helper()
}

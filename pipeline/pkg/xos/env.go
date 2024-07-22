// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// General environment variables.

package xos

import (
	"strings"
	"syscall"
)

const ExpansionStarterEnvironmentVariable = "$"

// Expand replaces ${var} or $var in the string based on the mapping function.
// For example, os.ExpandEnv(s) is equivalent to os.Expand(s, os.Getenv).
func Expand(s string, mapping func(string) string) string {
	return MultiExpand(s, []*Expander{{
		StartString: ExpansionStarterEnvironmentVariable,
		Fn:          mapping,
	}})
}

type Expander struct {
	StartString string
	Fn          func(string) string
}

// MultiExpand replaces ${var} or $var in the string based on the mapping function.
// For example, os.ExpandEnv(s) is equivalent to os.Expand(s, os.Getenv).
func MultiExpand(s string, expanders []*Expander) string {
	var buf []byte
	// ${} is all ASCII, so bytes are fine for this operation.
	i := 0
	j := 0
	for {
		var l int
		var expander *Expander

		expander, j, l = nextExpanderMatch(s, j, expanders)

		// if the last characters of the string match an expansion starter we need to break because otherwise we'll loop forever
		if expander == nil || j+l >= len(s) {
			break
		}

		if buf == nil {
			buf = make([]byte, 0, 2*len(s))
		}
		buf = append(buf, s[i:j]...)
		name, w := getShellName(s[j+l:])
		if name == "" && w > 0 {
			// Encountered invalid syntax; eat the
			// characters.
		} else if name == "" {
			// Valid syntax, but $ was not followed by a
			// name. Leave the dollar character untouched.
			buf = append(buf, s[j:j+l]...)
		} else {
			buf = append(buf, expander.Fn(name)...)
		}
		j += w + l
		i = j
	}
	if buf == nil {
		return s
	}
	return string(buf) + s[i:]
}

func nextExpanderMatch(s string, i int, expanders []*Expander) (expander *Expander, startIdx, length int) {
	for ; i < len(s); i++ {
		for _, expander := range expanders {
			if strings.HasPrefix(s[i:], expander.StartString) {
				return expander, i, len(expander.StartString)
			}
		}
	}

	return nil, -1, 0
}

// ExpandEnv replaces ${var} or $var in the string according to the values
// of the current environment variables. References to undefined
// variables are replaced by the empty string.
func ExpandEnv(s string) string {
	return Expand(s, Getenv)
}

// EscapedExpandEnv replaces ${var} or $var in the string according to the values
// of the current environment variables. References to undefined
// variables are replaced by the empty string.
func EscapedExpandEnv(s string) string {
	return Expand(s, EscapedGetEnv)
}

// isShellSpecialVar reports whether the character identifies a special
// shell variable such as $*.
func isShellSpecialVar(c uint8) bool {
	switch c {
	case '*', '#', '$', '@', '!', '?', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	}
	return false
}

// isAlphaNum reports whether the byte is an ASCII letter, number, or underscore
func isAlphaNum(c uint8) bool {
	return c == '_' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}

// getShellName returns the name that begins the string and the number of bytes
// consumed to extract it. If the name is enclosed in {}, it's part of a ${}
// expansion and two more bytes are needed than the length of the name.
func getShellName(s string) (string, int) {
	switch {
	case s[0] == '{':
		if len(s) > 2 && isShellSpecialVar(s[1]) && s[2] == '}' {
			return s[1:2], 3
		}
		// Scan to closing brace
		for i := 1; i < len(s); i++ {
			if s[i] == '}' {
				if i == 1 {
					return "", 2 // Bad syntax; eat "${}"
				}
				return s[1:i], i + 1
			}
		}
		return "", 1 // Bad syntax; eat "${"
	case isShellSpecialVar(s[0]):
		return s[0:1], 1
	}
	// Scan alphanumerics.
	var i int
	for i = 0; i < len(s) && isAlphaNum(s[i]); i++ {
	}
	return s[:i], i
}

// Getenv retrieves the value of the environment variable named by the key.
// It returns the value, which will be empty if the variable is not present.
// To distinguish between an empty value and an unset value, use LookupEnv.
func Getenv(key string) string {
	v, _ := syscall.Getenv(key)
	return v
}

func EscapedGetEnv(s string) string {
	// Done to escape dollar signs with another dollar sign
	if s == ExpansionStarterEnvironmentVariable {
		return ExpansionStarterEnvironmentVariable
	}

	return Getenv(s)
}

// LookupEnv retrieves the value of the environment variable named
// by the key. If the variable is present in the environment the
// value (which may be empty) is returned and the boolean is true.
// Otherwise the returned value will be empty and the boolean will
// be false.
func LookupEnv(key string) (string, bool) {
	return syscall.Getenv(key)
}

// Setenv sets the value of the environment variable named by the key.
// It returns an error, if any.
func Setenv(key, value string) error {
	err := syscall.Setenv(key, value)
	if err != nil {
		return NewSyscallError("setenv", err)
	}
	return nil
}

// Unsetenv unsets a single environment variable.
func Unsetenv(key string) error {
	return syscall.Unsetenv(key)
}

// Clearenv deletes all environment variables.
func Clearenv() {
	syscall.Clearenv()
}

// Environ returns a copy of strings representing the environment,
// in the form "key=value".
func Environ() []string {
	return syscall.Environ()
}

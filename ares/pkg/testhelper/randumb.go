package testhelper

import (
	"fmt"

	"github.com/Pallinder/go-randomdata"
)

var counter = 0

// This pretty dumb function prefixes a string with an auto incremented counter,
// to make random strings more unique.
func Randumb(s string) string {
	counter++

	return fmt.Sprintf("a%d-%s", counter, s)
}

func (h *Helper) Randumb(s string) string {
	return Randumb(s)
}

func RandomVersion() string {
	return fmt.Sprintf("v%d.%d.%d", randomdata.Number(1, 100), randomdata.Number(1, 100), randomdata.Number(1, 100))
}

func (h *Helper) RandomVersion() string {
	return RandomVersion()
}

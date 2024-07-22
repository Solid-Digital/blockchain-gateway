# A fork of the Go standard library's json encoder

## Why?
https://github.com/golang/go/issues/6213 was proposed in 2013 but was never accepted.

## Differences:
 - added support for the "inline" struct tag that forces the encoder/decoder to work as if the field was embedded into its parent struct

Examples:

```go
package main

import (
	"github.com/unchain/json"
)

type Child struct {
	Value string `json:"value"`
}

type Parent struct {
	Child *Child `json:",inline"`
}

type Parent2 struct {
	Child *Child `json:"child"`
}

func main() {
	bytes, err := json.Marshal(&Parent{
		Child: &Child{
			Value: "123",
		},
	})

	if err != nil || string(bytes) != `{"value":"123"}` {
		panic("panic")
	}
}

```
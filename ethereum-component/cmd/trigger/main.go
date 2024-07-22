package main

import (
	"bitbucket.org/unchain/ethereum2/pkg/trigger"
	"github.com/unchainio/interfaces/adapter"
)

var Version string

func NewTrigger() adapter.Trigger {
	return &trigger.Trigger{}
}

func main() {}

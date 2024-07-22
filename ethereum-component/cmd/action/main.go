package main

import (
	"bitbucket.org/unchain/ethereum2/pkg/action"
	"github.com/unchainio/interfaces/adapter"
)

var Version string

func NewAction() adapter.Action {
	return &action.Action{}
}

func main() {}

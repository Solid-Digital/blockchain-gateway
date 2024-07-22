package event

import "github.com/BurntSushi/toml"

type Config struct {
	ContractAddress string
	Name            string
	Filters         [][]toml.Primitive
}

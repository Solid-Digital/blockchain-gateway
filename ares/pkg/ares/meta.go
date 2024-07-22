package ares

import (
	"fmt"
	"runtime"

	"github.com/BurntSushi/toml"
)

// Metadata of the ares server
type Metadata struct {
	Name      string `toml:"name"`
	Version   string `toml:"version"`
	Branch    string `toml:"branch"`
	Builder   string `toml:"builder"`
	BuildDate string `toml:"buildDate"`
	GoVersion string `toml:"goVersion"`
}

// LoadMeta loads
func LoadMeta(name, metaTOML string) *Metadata {
	meta := new(Metadata)

	toml.Unmarshal([]byte(metaTOML), meta)

	if meta.Name == "" {
		meta.Name = name
	}

	meta.GoVersion = runtime.Version()

	return meta
}

func (m *Metadata) String() string {
	return fmt.Sprintf("%s:\n  version: %s\n  branch: %s\n  builder: %s\n  buildDate: %s\n  goVersion: %s\n", m.Name, m.Version, m.Branch, m.Builder, m.BuildDate, m.GoVersion)
}

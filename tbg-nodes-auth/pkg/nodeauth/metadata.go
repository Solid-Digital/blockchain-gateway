package nodeauth

import "fmt"

// Metadata of the ares server
type Metadata struct {
	Name      string `toml:"name"`
	Version   string `toml:"version"`
	Branch    string `toml:"branch"`
	Builder   string `toml:"builder"`
	BuildDate string `toml:"buildDate"`
	GoVersion string `toml:"goVersion"`
}

func (m *Metadata) String() string {
	return fmt.Sprintf("%s:\n  version: %s\n  branch: %s\n  builder: %s\n  buildDate: %s\n  goVersion: %s\n", m.Name, m.Version, m.Branch, m.Builder, m.BuildDate, m.GoVersion)
}

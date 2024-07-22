package factory

import "bitbucket.org/unchain/ares/pkg/ares"

func (f *Factory) Metadata() *ares.Metadata {
	return &ares.Metadata{
		Name:      "test_ares",
		Version:   "0.1",
		BuildDate: "01-01-2001",
	}
}

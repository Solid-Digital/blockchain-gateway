package pipeline

import (
	"io/ioutil"

	"github.com/unchainio/pkg/errors"
)

type Config struct {
	Trigger ComponentConfig
	Actions []ComponentConfig
}

type ComponentConfig struct {
	Name          string
	Path          string
	Config        string
	ConfigPath    string
	MessageConfig map[string]interface{}
}

func (a *ComponentConfig) ConfigBytes() ([]byte, error) {
	if a.ConfigPath != "" {
		configBytes, err := ioutil.ReadFile(a.ConfigPath)
		if err != nil {
			return nil, errors.Wrapf(err, "could not find config file %s ", a.ConfigPath)
		}
		return configBytes, nil
	} else {
		return []byte(a.Config), nil
	}
}

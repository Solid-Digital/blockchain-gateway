package factory

import (
	"bytes"
	"io/ioutil"

	"bitbucket.org/unchain/ares/pkg/component"
)

func (f *Factory) File(filePath, fileName string) {
	f.suite.Require().NotNil(f.ares)

	b, err := ioutil.ReadFile(filePath)

	f.suite.Require().NoError(err)

	reader := bytes.NewReader(b)

	tar, err := component.TarGz(fileName, reader)

	f.suite.Require().NoError(err)

	err = f.ares.FileStore.PutObject(fileName, tar)

	f.suite.Require().NoError(err)
}

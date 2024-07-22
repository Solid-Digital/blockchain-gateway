package component

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"time"

	"github.com/unchainio/pkg/errors"
)

func TarGz(fileName string, reader io.Reader) (io.Reader, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	buf := new(bytes.Buffer)

	gzw := gzip.NewWriter(buf)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	mode := int64(0600)

	header := &tar.Header{
		Name:       fileName,
		ChangeTime: time.Now().UTC(),
		Mode:       mode,
		Size:       int64(len(data)),
	}

	if err := tw.WriteHeader(header); err != nil {
		return nil, errors.Wrap(err, "failed to write tar header")
	}

	if _, err := tw.Write(data); err != nil {
		return nil, errors.Wrap(err, "failed to write tar file")
	}

	return buf, nil
}

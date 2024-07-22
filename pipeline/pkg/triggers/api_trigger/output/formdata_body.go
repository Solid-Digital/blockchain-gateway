package output

import (
	"bytes"
	"io"
	"net/http"
)

func formdataBody(req *http.Request) (output []byte, err error) {
	var buf bytes.Buffer
	file, _, err := req.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(&buf, file)

	return buf.Bytes(), nil
}


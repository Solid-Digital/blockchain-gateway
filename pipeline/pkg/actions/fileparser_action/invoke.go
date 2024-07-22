package fileparser_action

import (
	"github.com/unchain/pipeline/pkg/domain"
	"github.com/unchainio/pkg/errors"
	"strings"
)

const (
	FileType = "filetype"
	File = "file"
	Header = "header"
	Delimiter = "delimiter"
)

func Invoke(stub domain.Stub, input map[string]interface{}) (output map[string]interface{}, err error) {
	filetype, ok := input[FileType].(string)
	if (!ok) {
		return nil, errors.New("could not cast filetype to string")
	}
	rawFile, ok := input[File].([]byte)
	if (!ok) {
		return nil, errors.New("could not cast file to []byte")
	}


	switch strings.ToLower(filetype) {
	case "csv":
		// fallback config for CSV
		header := false
		if input[Header] != nil {
			header, ok = input[Header].(bool)
			if (!ok) {
				return nil, errors.New("could not cast header to bool")
			}
		}
		delimiter := ','
		if input[Delimiter] != nil {
			delimiter, ok = input[Delimiter].(rune)
			if (!ok) {
				return nil, errors.New("could not cast delimiter to rune")
			}
		}
		return ParseCsv(stub, rawFile, header, delimiter)
	case "json":
		return ParseJson(stub, rawFile)
	default:
		return nil, errors.New("unknown filetype")
	}
}
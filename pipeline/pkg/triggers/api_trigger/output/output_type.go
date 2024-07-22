package output

import (
	"strings"
)

const (
	JSON = iota
	XML
	FORM
	CSV
	TEXT
	FORMDATA
	UNKNOWN
)

func outputType(contentType string) int {
	if strings.Contains(contentType, "multipart/form-data") {
		return FORMDATA
	}

	switch contentType {
	case "application/json":
		return JSON
	case "application/xml", "text/xml":
		return XML
	case "application/x-www-form-urlencoded":
		return FORM
	default:
		if len(contentType) > 5 && contentType[:5] == "text/" {
			return TEXT
		}

		return UNKNOWN
	}
}

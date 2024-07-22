package fileparser_action

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/unchain/pipeline/pkg/domain"
	"strings"
)

func ParseCsv(stub domain.Stub, rawFile []byte, header bool, delimiter rune) (map[string]interface{}, error) {
	r := csv.NewReader(bytes.NewReader(rawFile))
	r.Comma = delimiter // use ; delimiter instead of ,

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	// store records as array of map[string]interface{}
	var headerKeys []string
	var recordsArray []map[string]interface{}
	for i, record := range records {
		if header && i == 0 {
			headerKeys = getHeaderKeys(record)
		} else {
			recordMap := createMap(headerKeys, record)
			recordsArray = append(recordsArray, recordMap)
		}
	}
	stub.Debugf("records: %v", recordsArray)

	return map[string]interface{}{
		"messages": recordsArray,
	}, nil
}

func getHeaderKeys(record []string) []string {
	headerKeys := make([]string, len(record))
	for i, s := range record {
		headerKeys[i] = strings.Join(strings.Fields(s), "")
	}
	return headerKeys
}

func createMap(headerKeys []string, record []string) map[string]interface{} {
	output := make(map[string]interface{})
	for i, v := range record {
		if headerKeys != nil {
			output[headerKeys[i]] = v
		} else {
			output[fmt.Sprintf("col-%v", i)] = v
		}
	}
	return output
}
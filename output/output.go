package output

import (
	"encoding/json"
	"fmt"
	"os"
)

type Field struct {
	Field             string `json:"field"`
	File              string `json:"file"`
	Line              int    `json:"line"`
	DeprecationReason string `json:"reason"`
}

type Data []Field

func CompareFiles(fileA string, fileB string) (Data, error) {
	aContent, err := os.ReadFile(fileA)
	if err != nil {
		return nil, err
	}

	bContent, err := os.ReadFile(fileB)
	if err != nil {
		return nil, err
	}

	aData := Data{}
	err = json.Unmarshal(aContent, &aData)
	if err != nil {
		return nil, fmt.Errorf("could not json decode %s: %w", fileA, err)
	}

	bData := Data{}
	err = json.Unmarshal(bContent, &bData)
	if err != nil {
		return nil, fmt.Errorf("could not json decode %s: %w", fileB, err)
	}

	result := Data{}
	// This is a dumb loop that could be optimized
	for _, bField := range bData {
		found := false
		for _, aField := range aData {
			if bField.Field == aField.Field {
				found = true
				break
			}
		}
		if !(found) {
			result = append(result, bField)
		}
	}

	return result, nil
}

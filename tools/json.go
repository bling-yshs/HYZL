package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func GetValueFromJSONFile(filePath, key string) (interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	var jsonData map[string]interface{}
	if err = json.Unmarshal(byteValue, &jsonData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	value, ok := jsonData[key]
	if !ok {
		return nil, fmt.Errorf("key '%s' not found in JSON file", key)
	}

	return value, nil
}

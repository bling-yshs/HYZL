package tools

import (
	"encoding/json"
	"errors"
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

	var jsonData interface{}
	if err = json.Unmarshal(byteValue, &jsonData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	value, err := getJsonValue(jsonData, key)
	if err != nil {
		return nil, fmt.Errorf("key '%s' not found in JSON file", key)
	}

	return value, nil
}

func CheckKeyInJSONFile(filePath, key string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return false, fmt.Errorf("failed to read file: %v", err)
	}

	var jsonData interface{}
	if err = json.Unmarshal(byteValue, &jsonData); err != nil {
		return false, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	_, err = getJsonValue(jsonData, key)
	if err == nil {
		return true, nil
	} else {
		return false, nil
	}

}

func getJsonValue(data interface{}, key string) (interface{}, error) {
	switch v := data.(type) {
	case map[string]interface{}:
		for k, v := range v {
			if k == key {
				return v, nil
			}
			value, err := getJsonValue(v, key)
			if err == nil {
				return value, nil
			}
		}
	case []interface{}:
		for _, v := range v {
			value, err := getJsonValue(v, key)
			if err == nil {
				return value, nil
			}
		}
	}
	return nil, errors.New("key not found")
}

func UpdateValueInJSONFile(filePath, key, subKey string, newValue interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	var jsonData map[string]interface{}
	if err = json.Unmarshal(byteValue, &jsonData); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	subObj, ok := jsonData[key]
	if !ok {
		return fmt.Errorf("key '%s' not found in JSON file", key)
	}

	subData, ok := subObj.(map[string]interface{})
	if !ok {
		return fmt.Errorf("key '%s' is not a JSON object", key)
	}

	subData[subKey] = newValue
	jsonData[key] = subData

	updatedJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated JSON: %v", err)
	}

	err = os.WriteFile(filePath, updatedJSON, 0644)
	if err != nil {
		return fmt.Errorf("failed to write updated JSON to file: %v", err)
	}

	return nil
}

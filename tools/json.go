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

func CheckKeyInJSONFile(filePath, key, subKey string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return false, fmt.Errorf("failed to read file: %v", err)
	}

	var jsonData map[string]interface{}
	if err = json.Unmarshal(byteValue, &jsonData); err != nil {
		return false, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	if subKey != "" {
		subObj, ok := jsonData[key]
		if !ok {
			return false, fmt.Errorf("key '%s' not found in JSON file", key)
		}

		subData, ok := subObj.(map[string]interface{})
		if !ok {
			return false, fmt.Errorf("key '%s' is not a JSON object", key)
		}

		_, ok = subData[subKey]
		return ok, nil
	}

	_, ok := jsonData[key]

	return ok, nil
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

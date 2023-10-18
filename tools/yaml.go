package tools

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

func GetValueFromYAMLFile(filePath, key string) (interface{}, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %v", err)
	}

	data := make(map[string]interface{})
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		return "", fmt.Errorf("解析YAML失败: %v", err)
	}

	value, ok := data[key]
	if !ok {
		return "", fmt.Errorf("键 '%s' 不存在", key)
	}

	return value, nil
}

func UpdateOrAppendToYaml(filePath string, key string, value interface{}) error {
	// 读取 YAML 文件
	yamlContent, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// 解析 YAML
	node := &yaml.Node{}
	err = yaml.Unmarshal(yamlContent, node)
	if err != nil {
		return err
	}

	// 检查是否存在同名的键
	found := false
	for i := 0; i < len(node.Content[0].Content); i += 2 {
		if node.Content[0].Content[i].Value == key {
			// 修改键的值
			if node.Content[0].Content[i+1].Tag == "!!null" {
				// 如果是 null 类型，替换整个节点
				node.Content[0].Content[i+1] = &yaml.Node{
					Kind:  yaml.ScalarNode,
					Value: fmt.Sprintf("%v", value),
				}
			} else {
				node.Content[0].Content[i+1].Value = fmt.Sprintf("%v", value)
			}
			found = true
			break
		}
	}

	if !found {
		// 在末尾添加新键和值
		newKeyValue := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: key,
		}
		newValue := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: fmt.Sprintf("%v", value),
		}

		node.Content[0].Content = append(node.Content[0].Content, newKeyValue, newValue)
	}

	// 将修改后的 YAML 写回到 filePath
	outputFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	encoder := yaml.NewEncoder(outputFile)
	encoder.SetIndent(2)
	err = encoder.Encode(node)
	if err != nil {
		return err
	}

	return nil
}

func UpdateValueYAML(filePath string, key string, value interface{}) error {
	// 读取 YAML 文件
	yamlContent, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// 解析 YAML
	node := &yaml.Node{}
	err = yaml.Unmarshal(yamlContent, node)
	if err != nil {
		return err
	}

	// 检查是否存在同名的键
	found := false
	for i := 0; i < len(node.Content[0].Content); i += 2 {
		if node.Content[0].Content[i].Value == key {
			// 修改键的值
			node.Content[0].Content[i+1].Value = fmt.Sprintf("%v", value)
			found = true
			break
		}
	}
	if found {
		// 将修改后的 YAML 写回到 filePath
		outputFile, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer outputFile.Close()
		encoder := yaml.NewEncoder(outputFile)
		encoder.SetIndent(2)
		err = encoder.Encode(node)
		if err != nil {
			return err
		}
	}

	return nil
}

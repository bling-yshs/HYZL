package yaml_utils

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"os"
)

func GetValueFromYAMLFile(filePath, key string) (interface{}, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return "", errors.Wrap(err, "错误描述：读取文件失败")
	}

	data := make(map[string]interface{})
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		return "", errors.Wrap(err, "错误描述：解析YAML失败")
	}

	value, ok := data[key]
	if !ok {
		return "", errors.Wrapf(err, "键 '%s' 不存在", key)
	}

	return value, nil
}

func UpdateOrAppendToYaml(filePath string, key string, value interface{}) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	node := &yaml.Node{}
	err = yaml.Unmarshal(file, node)
	if err != nil {
		return err
	}
	var found bool = false
	for index, item := range node.Content[0].Content {
		if item.Value == key {
			found = true
			// 如果value的类型是nil，则将tag设置为!!null，value设置为空字符串，style设置为0
			if value == nil {
				node.Content[0].Content[index+1].Tag = "!!null"
				node.Content[0].Content[index+1].Value = ""
				node.Content[0].Content[index+1].Style = 0
				break
			}
			// 如果值是string，就将tag设置为!!str，value设置为传入的值
			if strValue, ok := value.(string); ok {
				node.Content[0].Content[index+1].Tag = "!!str"
				node.Content[0].Content[index+1].Value = strValue
				break
			}
		}
	}
	if !found {
		// 如果value的类型是nil，则将tag设置为!!null，value设置为空字符串，style设置为0
		if value == nil {
			node.Content[0].Content = append(node.Content[0].Content, &yaml.Node{
				Kind:  yaml.ScalarNode,
				Tag:   "!!str",
				Value: key,
			})
			node.Content[0].Content = append(node.Content[0].Content, &yaml.Node{
				Kind:  yaml.ScalarNode,
				Tag:   "!!null",
				Value: "",
				Style: 0,
			})
		}
	}
	file, err = yaml.Marshal(node)
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, file, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func UpdateValueYAML(filePath string, key string, value interface{}) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	node := &yaml.Node{}
	err = yaml.Unmarshal(file, node)
	if err != nil {
		return err
	}
	for index, item := range node.Content[0].Content {
		if item.Value == key {
			// 如果value的类型是nil，则将tag设置为!!null，value设置为空字符串，style设置为0
			if value == nil {
				node.Content[0].Content[index+1].Tag = "!!null"
				node.Content[0].Content[index+1].Value = ""
				node.Content[0].Content[index+1].Style = 0
				break
			}
			// 如果值是string，就将tag设置为!!str，value设置为传入的值
			if strValue, ok := value.(string); ok {
				node.Content[0].Content[index+1].Style = yaml.SingleQuotedStyle
				node.Content[0].Content[index+1].Tag = "!!str"
				node.Content[0].Content[index+1].Value = strValue
				break
			}
			// 如果值是int，就将tag设置为!!int，value设置为传入的值
			if uintValue, ok := value.(uint32); ok {
				node.Content[0].Content[index+1].Tag = "!!int"
				node.Content[0].Content[index+1].Value = fmt.Sprintf("%d", uintValue)
				break
			}
		}
	}
	file, err = yaml.Marshal(node)
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, file, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func DeleteKey(filePath string, key string) error {
	// 读取文件
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	// 解析yaml
	node := &yaml.Node{}
	err = yaml.Unmarshal(file, node)
	if err != nil {
		return err
	}
	// 遍历yaml的内容
	for index, item := range node.Content[0].Content {
		// 如果找到了key
		if item.Value == key {
			// 删除key和value
			node.Content[0].Content = append(node.Content[0].Content[:index], node.Content[0].Content[index+2:]...)
			break
		}
	}
	// 重新序列化yaml
	file, err = yaml.Marshal(node)
	if err != nil {
		return err
	}
	// 写入文件
	err = os.WriteFile(filePath, file, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

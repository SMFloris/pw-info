package pipewire

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

type PipewireNodeProperty struct {
	Key   string
	Value string
}

type PipewireNode struct {
	Id         int
	Properties []PipewireNodeProperty
	Type       string
}

func ParseOutput(output string) []PipewireNode {
	var lines []string = strings.Split(output, "\n")

	var result []PipewireNode

	var currentNode PipewireNode
	currentNode.Properties = []PipewireNodeProperty{}
	var inProperties bool = false
	for _, line := range lines {
		trimmedLine := strings.ReplaceAll(line, " ", "")

		if inProperties && !strings.Contains(trimmedLine, "\t\t") {
			inProperties = false
		}

		// if id is encountered then it means we are starting a new property
		if strings.Contains(trimmedLine, "\tid:") {
			inProperties = false
			intVar, err := strconv.Atoi(strings.TrimPrefix(trimmedLine, "\tid:"))
			if err != nil {
				log.Fatal(err)
			}
			if currentNode.Id > 0 {
				result = append(result, currentNode)
			}
			currentNode.Id = intVar
			currentNode.Properties = []PipewireNodeProperty{}
			continue
		}

		if strings.Contains(trimmedLine, "\ttype:") {
			value := strings.TrimPrefix(trimmedLine, "\ttype:")
			currentNode.Type = value
		}

		if strings.Contains(trimmedLine, "\tproperties:") {
			inProperties = true
			continue
		}

		if !inProperties {
			continue
		}

		// we are inside properties
		var key string
		var value string
		var found bool
		key, value, found = strings.Cut(trimmedLine, "=")
		if found {
			key := strings.ReplaceAll(key, "\"", "")
			key = strings.ReplaceAll(key, "\t", "")
			value := strings.ReplaceAll(value, "\"", "")
			value = strings.ReplaceAll(value, "\t", "")
			var property PipewireNodeProperty
			property.Key = key
			property.Value = value
			currentNode.Properties = append(currentNode.Properties, property)
		}
	}
	result = append(result, currentNode)

	return result
}

func GetProperty(nodes []PipewireNode, id int, name string) (PipewireNodeProperty, error) {
	properties, err := ListProperties(nodes, id)
	if err != nil {
		return PipewireNodeProperty{}, err
	}

	for _, property := range properties {
		if property.Key == name {
			return property, nil
		}
	}

	return PipewireNodeProperty{}, errors.New("No property with that name found")
}

func ListProperties(nodes []PipewireNode, id int) ([]PipewireNodeProperty, error) {
	for _, node := range nodes {
		if node.Id != id {
			continue
		}
		return node.Properties, nil
	}

	return []PipewireNodeProperty{}, errors.New("No node found with that id")
}

func FindNode(nodes []PipewireNode, key string, value string) (PipewireNode, error) {
	for _, node := range nodes {
		for _, property := range node.Properties {
			if property.Key == key && property.Value == value {
				return node, nil
			}
		}
	}

	return PipewireNode{}, errors.New("No node found")
}

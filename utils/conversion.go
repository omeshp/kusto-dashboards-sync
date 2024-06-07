package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

func ConvertYAMLtoJSON(yamlFile string, jsonFile string) error {
	jsonData, err := ConvertYAMLToJSON(yamlFile)
	if err != nil {
		return fmt.Errorf("error marshalling JSON data: %v", err)
	}

	// Write JSON to output.json file
	err = ioutil.WriteFile(jsonFile, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("Error writing JSON to file: %v\n", err)
	}

	fmt.Printf("JSON data written to %s\n", yamlFile)
	return nil
}

// ConvertYAMLToJSON converts a YAML file to JSON format
func ConvertYAMLToJSON(yamlFile string) ([]byte, error) {
	// Read the YAML file
	yamlData, err := os.ReadFile(yamlFile)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %v", err)
	}

	// Unmarshal the YAML into an interface{}
	var yamlObj interface{}
	err = yaml.Unmarshal(yamlData, &yamlObj)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML: %v", err)
	}

	// Marshal the interface{} to JSON
	jsonData, err := json.MarshalIndent(yamlObj, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshaling to JSON: %v", err)
	}

	return jsonData, nil
}

func ConvertJSONtoYAML(jsonFile, yamlFile string) error {
	// Convert map to YAML
	yamlData, err := ConvertJSONToYAML(jsonFile)
	if err != nil {
		return fmt.Errorf("error marshalling YAML data: %v", err)
	}

	// Write YAML data to a file
	err = ioutil.WriteFile(yamlFile, yamlData, 0644)
	if err != nil {
		return fmt.Errorf("error writing YAML data to file: %v", err)
	}

	fmt.Printf("YAML data written to %s\n", yamlFile)
	return nil
}

// ConvertJSONToYAML converts a JSON file to YAML format
func ConvertJSONToYAML(jsonFile string) ([]byte, error) {
	// Read the JSON file
	jsonData, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON file: %v", err)
	}

	// Unmarshal the JSON into an interface{}
	var jsonObj interface{}
	err = json.Unmarshal(jsonData, &jsonObj)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	// Marshal the interface{} to YAML
	yamlData, err := yaml.Marshal(jsonObj)
	if err != nil {
		return nil, fmt.Errorf("error marshaling to YAML: %v", err)
	}

	return yamlData, nil
}

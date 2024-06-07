package utils

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

// include function reads the file content and returns it as a string
func include(filename string) (string, error) {
	filepath := fmt.Sprintf("queries/%s", filename)
	// replace all escaped single quotes with single quotes
	filepath = strings.ReplaceAll(filepath, "''", "'")
	content, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	value := string(content)

	value = strings.ReplaceAll(value, "\n", "\\n")
	value = strings.ReplaceAll(value, "'", "''")
	value = "'" + value + "'"

	return value, nil
}

// ProcessTemplate processes the YAML template and writes the output to a file
func ProcessTemplate(templatePath, outputPath string) error {
	// Read the template file
	tmplContent, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("error reading template file: %w", err)
	}

	// Create a new template and register the include function
	tmpl, err := template.New("yamlTemplate").Funcs(template.FuncMap{
		"include": include,
	}).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	// Create the output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer outFile.Close()

	// Execute the template with an empty data context
	if err = tmpl.Execute(outFile, nil); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil
}

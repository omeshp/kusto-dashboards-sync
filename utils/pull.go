package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/omeshp/kusto-dashboards-sync/models"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"regexp"
	"strings"
)

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

func ConvertRawDashboardToConcrete(rawDashboard *interface{}) (*models.Dashboard, error) {
	// Marshal the raw dashboard into JSON
	rawDashboardJSON, err := json.Marshal(rawDashboard)
	if err != nil {
		return nil, fmt.Errorf("error marshalling raw dashboard: %v", err)
	}

	// Unmarshal the JSON into a models.Dashboard
	var dashboard models.Dashboard
	if err := json.Unmarshal(rawDashboardJSON, &dashboard); err != nil {
		return nil, fmt.Errorf("error unmarshalling dashboard data: %v", err)
	}

	return &dashboard, nil
}

func PersistDashboardData(dashboardRaw *interface{}, masterDashboard *models.Dashboard, outputYamlPath string) error {
	if _, err := os.Stat("queries"); !os.IsNotExist(err) {
		// If the directory exists, delete it
		err := os.RemoveAll("queries")
		if err != nil {
			log.Fatalf("Failed to clean directory: %v", err)
		}
	}

	// Create queries directory
	err := os.MkdirAll("queries", 0755)
	if err != nil {
		return fmt.Errorf("error creating queries directory: %v", err)
	}

	//// Marshal YAML data
	//yamlBytes, err := yaml.Marshal(dashboardRaw)
	//if err != nil {
	//	fmt.Printf("Error marshaling YAML: %v\n", err)
	//	return fmt.Errorf("error reading yaml bytes: %v", err)
	//}
	//
	//var data interface{}
	//err = yaml.Unmarshal(yamlBytes, &data)
	//if err != nil {
	//	return fmt.Errorf("Error unmarshaling YAML: %v\n", err)
	//}
	// The YAML data is now in a nested map structure
	dataMap := (*dashboardRaw).(map[string]interface{})

	dashboard, err := ConvertRawDashboardToConcrete(dashboardRaw)
	if err != nil {
		fmt.Println("Error:", err)
		return fmt.Errorf("error converting to Dashboard struct: %v", err)
	}

	for tileIndex, tile := range dashboard.Tiles {
		var query models.Query
		currentPageName := ""
		for _, page := range dashboard.Pages {
			if tile.PageId == page.Id {
				currentPageName = page.Name
			}
		}

		filename := currentPageName + "_" + tile.Title
		if tile.VisualType == "markdownCard" {
			filename = strings.ReplaceAll(filename, " ", "_") + ".md"
		} else {
			filename = strings.ReplaceAll(filename, " ", "_") + ".kql"
		}

		if tile.QueryRef.QueryId != "" {
			// Find the query by QueryRef
			for queryIndex, q := range dashboard.Queries {
				if q.Id == tile.QueryRef.QueryId {
					query = q
					dataMap["queries"].([]interface{})[queryIndex].(map[string]interface{})["text"] = "{{ include " + "\"" + filename + "\"}}"
					break
				}
			}
		} else if tile.Query.Text != "" {
			// Use the Query directly from the tile
			query = tile.Query
			var s interface{} = `{{ include "` + filename + `"}}`
			dataMap["tiles"].([]interface{})[tileIndex].(map[string]interface{})["query"].(map[string]interface{})["text"] = s
		}

		if query.Text != "" {
			filepath := fmt.Sprintf("queries/%s", filename)
			err := os.WriteFile(filepath, []byte(query.Text), 0644)
			if err != nil {
				return fmt.Errorf("error writing data to file %s: %v", filename, err)
			}
			fmt.Printf("Query saved to file: %s\n", filepath)
		}
	}

	// retain id, title, etag from master dashboard
	dataMap["id"] = masterDashboard.Id
	dataMap["title"] = masterDashboard.Title
	dataMap["eTag"] = masterDashboard.ETag

	// Marshal the data back into a YAML string
	newYamlBytes, err := yaml.Marshal(&dataMap)
	if err != nil {
		return fmt.Errorf("Error marshaling YAML: %v\n", err)
	}

	yamlData := string(newYamlBytes)
	re := regexp.MustCompile(`'{{ include "(.*?)"}}'`)
	yamlData = re.ReplaceAllString(yamlData, `{{ include "$1"}}`)

	err = os.WriteFile(outputYamlPath, []byte(yamlData), 0644)
	if err != nil {
		fmt.Println("Error:", err)
		// handle the error
	}

	fmt.Printf("Saved dashboard template to: %s\n", outputYamlPath)

	return nil
}

// replaceQueryText replaces the query text with include filename in jsonString
func replaceQueryText(jsonString, queryText, filename string) string {
	return strings.ReplaceAll(jsonString, "\""+queryText+"\"", "{{ include "+"\""+filename+"\"}}")
}

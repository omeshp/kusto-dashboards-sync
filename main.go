package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
	"kusto-dashboards-sync/dataexplorer"
	"kusto-dashboards-sync/models"
	"kusto-dashboards-sync/utils"
	"log"
	"os"
	"strings"
)

const Dashboard_Template_Path = "dashboard.yml"
const Dashboard_Output_Path = "bin/dashboard_processed.yml"
const Dashboard_JSON_Output_Path = "bin/dashboard.json"

func main() {

	// Define the command-line arguments
	var command string
	var dashboardID = ""

	// Customize the usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Println("  pull: Pull the data for dashboard set in config.yml")
		fmt.Println("  push: Push the data to dashboard set in config.yml")
		fmt.Println("  pull [dashboard id]")
		fmt.Println("  push [dashboard id]")
	}

	// Parse the command-line arguments
	flag.Parse()

	// Check for remaining arguments (non-flag arguments)
	if flag.NArg() > 0 {
		command = flag.Arg(0)
	}

	if flag.NArg() > 1 {
		dashboardID = flag.Arg(1)
	}

	fmt.Printf("Command: %s\n", command)
	fmt.Printf("Dashboard ID: %s\n", dashboardID)

	// If no flags are provided, print the usage message and exit
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	} else {
		command = flag.Arg(0)
	}

	config, err := getDashboardConfig()
	if err != nil {
		log.Fatalf("Error loading dashboard config from config.yml file")
	}

	// Load environment variables from .env file
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	if dashboardID != "" {
		fmt.Printf("Overriding with Dashboard ID: %s\n", dashboardID)
		config.DashboardID = dashboardID
	}

	// Get the access token from the environment variables
	accessToken := os.Getenv("ACCESS_TOKEN")
	if accessToken == "" {
		log.Fatalf("ACCESS_TOKEN not set in .env file")
	}

	accessToken = "Bearer " + accessToken

	// Create queries directory
	err = os.MkdirAll("bin", 0755)
	if err != nil {
		fmt.Errorf("error creating queries directory: %v", err)
	}

	if command == "pull" {
		PullDashboard(config.DashboardID, accessToken, err)
	}

	if command == "push" {
		PushDashboard(err, accessToken, config.DashboardID)
	}

}

type Config struct {
	DashboardID string `yaml:"dashboard_id"`
}

func getDashboardConfig() (*Config, error) {
	data, err := os.ReadFile("config.yml")
	if err != nil {
		return &Config{}, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return &Config{}, err
	}

	return &cfg, nil
}

func PushDashboard(err error, accessToken string, dashboardId string) {
	err = utils.ProcessTemplate(Dashboard_Template_Path, Dashboard_Output_Path)
	if err != nil {
		fmt.Println("Failed to process template file: %v\n", err)
		return
	}

	fmt.Printf("Succeeded in processing template file: %s, with output: %s\n", Dashboard_Template_Path, Dashboard_Output_Path)

	// Convert YAML to JSON
	jsonData, err := utils.ConvertYAMLToJSON(Dashboard_Output_Path)
	if err != nil {
		log.Fatalf("Failed to convert YAML to JSON: %v", err)
	}

	jsonString := string(jsonData)
	jsonString = strings.ReplaceAll(jsonString, "\\\\n", "\\n")

	// Write JSON to output.json file
	err = os.WriteFile(Dashboard_JSON_Output_Path, []byte(jsonString), 0644)
	if err != nil {
		fmt.Printf("Error writing JSON to file: %v\n", err)
		return
	}

	fmt.Printf("Dashboard json written to %s\n", Dashboard_JSON_Output_Path)
	// Unmarshal the response body into a Dashboard struct
	var dashboard interface{}
	if err := json.Unmarshal([]byte(jsonString), &dashboard); err != nil {
		log.Fatalf("Failed to convert YAML to JSON: %v", err)
	}

	dataExplorerClient := dataexplorer.NewDataExplorerClient("https://dashboards.kusto.windows.net/dashboards/", accessToken)

	currentDashboard, err := dataExplorerClient.GetDashboard(dashboardId)
	if err != nil {
		log.Fatalf("error updating dashboard: %v", err)
	}

	dashboardMap := dashboard.(map[string]interface{})
	dashboardMap["eTag"] = currentDashboard.ETag

	// Call the function to update the dashboard
	err = dataExplorerClient.UpdateDashboardRaw(dashboardId, &dashboard)
	if err != nil {
		log.Fatalf("error updating dashboard: %v", err)
	}
}

func PullDashboard(dashboardID string, accessToken string, err error) {
	dataExplorerClient := dataexplorer.NewDataExplorerClient("https://dashboards.kusto.windows.net/dashboards/", accessToken)

	// Get dashboard
	rawDashboard, err := dataExplorerClient.GetDashboardRaw(dashboardID)
	if err != nil {
		fmt.Printf("Error retrieving dashboard: %v\n", err)
		return
	}

	// Convert raw dashboard to concrete one
	dashboard, err := ConvertRawDashboardToConcrete(rawDashboard)
	if err != nil {
		fmt.Printf("Error converting raw dashboard to concrete one: %v\n", err)
		return
	}

	fmt.Printf("Retrieved Dashboard ID: %s, Title: %s\n", dashboardID, dashboard.Title)

	// Save queries to files
	err = utils.PersistDashboardData(rawDashboard, Dashboard_Template_Path)
	if err != nil {
		log.Fatalf("error saving queries to files: %v", err)
	}
}

// ConvertRawDashboardToConcrete takes a raw dashboard and converts it into a models.Dashboard
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

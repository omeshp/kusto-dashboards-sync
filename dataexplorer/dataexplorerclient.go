package dataexplorer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"kusto-dashboards-sync/models"
	"log"
	"net/http"
)

// DataExplorerClient represents a Data Explorer client
type DataExplorerClient struct {
	Client  *http.Client
	BaseURL string
}

// NewDataExplorerClient creates a new instance of DataExplorerClient
func NewDataExplorerClient(baseURL string, accessToken string) *DataExplorerClient {
	client := &http.Client{}
	client.Transport = &Transport{
		AccessToken: accessToken,
	}

	return &DataExplorerClient{
		Client:  client,
		BaseURL: baseURL,
	}
}

// Transport is a custom RoundTripper that adds Authorization header to each request
type Transport struct {
	AccessToken string
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", t.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	return http.DefaultTransport.RoundTrip(req)
}

// GetDashboard fetches the dashboard using the provided ID
func (dec *DataExplorerClient) GetDashboard(dashboardID string) (*models.Dashboard, error) {
	url := fmt.Sprintf("%s/%s", dec.BaseURL, dashboardID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := dec.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return nil, fmt.Errorf("error response from server: %s", bodyString)
	}

	var dashboard models.Dashboard
	err = json.NewDecoder(resp.Body).Decode(&dashboard)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &dashboard, nil
}

// GetDashboardRaw retrieves a dashboard using a GET call with the provided HTTP client and returns the dashboard data or an error
func (dec *DataExplorerClient) GetDashboardRaw(dashboardID string) (*interface{}, error) {
	// Make a GET request to retrieve the dashboard
	resp, err := dec.Client.Get(fmt.Sprintf("%s/%s", dec.BaseURL, dashboardID))
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error retrieving dashboard: status code %d", resp.StatusCode)
	}

	// Unmarshal the response body into a Dashboard struct
	var dashboard interface{}
	if err := json.Unmarshal(body, &dashboard); err != nil {
		return nil, fmt.Errorf("error unmarshalling dashboard data: %v", err)
	}

	return &dashboard, nil
}

// UploadDashboard uploads the dashboard using the provided ID
func (dec *DataExplorerClient) UploadDashboard(dashboardID string, dashboard *models.Dashboard) error {
	url := fmt.Sprintf("%s/%s", dec.BaseURL, dashboardID)
	dashboardJSON, err := json.MarshalIndent(dashboard, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling dashboard: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(dashboardJSON))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	resp, err := dec.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return fmt.Errorf("error response from server: %s", bodyString)
	}

	log.Printf("Successfully uploaded dashboard with ID %s. HTTP Status: %s", dashboardID, resp.Status)
	return nil
}

// UpdateDashboardRaw uploads a dashboard using a POST call with the provided HTTP client and returns an error if any
func (dec *DataExplorerClient) UpdateDashboardRaw(dashboardId string, dashboard *interface{}) error {
	// Marshal the dashboard data into JSON
	payload, err := json.Marshal(dashboard)
	if err != nil {
		return fmt.Errorf("error marshalling dashboard data: %v", err)
	}

	// Create a PUT request to upload the dashboard
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("https://dashboards.kusto.windows.net/dashboards/%s", dashboardId), bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("error creating PUT request: %v", err)
	}

	// Send the PUT request
	resp, err := dec.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error making PUT request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error uploading dashboard: status code %d, response body: %s", resp.StatusCode, string(body))
	}

	fmt.Printf("Dashboard updated successfully\n")

	return nil
}

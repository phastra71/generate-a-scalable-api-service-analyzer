package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// APIEndpoint represents an API endpoint
type APIEndpoint struct {
	Method  string `json:"method"`
	Path    string `json:"path"`
	Query   string `json:"query"`
	Headers string `json:"headers"`
	Body    string `json:"body"`
}

// APIService represents an API service
type APIService struct {
	Name     string        `json:"name"`
	Endpoints []APIEndpoint `json:"endpoints"`
}

// APIAnalyzer analyzes an API service
type APIAnalyzer struct {
	log *logrus.Logger
}

// NewAPIAnalyzer returns a new APIAnalyzer instance
func NewAPIAnalyzer() *APIAnalyzer {
	return &APIAnalyzer{
		log: logrus.New(),
	}
}

// Analyze analyzes an API service
func (a *APIAnalyzer) Analyze(service *APIService) {
	a.log.Infof("Analyzing API service '%s'", service.Name)

	for _, endpoint := range service.Endpoints {
		a.log.Debugf("Analyzing endpoint '%s %s'", endpoint.Method, endpoint.Path)

		// Send request to endpoint
		req, err := http.NewRequest(endpoint.Method, endpoint.Path, strings.NewReader(endpoint.Body))
		if err != nil {
			a.log.Errorf("Error creating request: %v", err)
			continue
		}

		// Set query parameters
		req.URL.RawQuery = endpoint.Query

		// Set headers
		headers := strings.Split(endpoint.Headers, ",")
		for _, header := range headers {
			keyValue := strings.Split(header, ":")
			req.Header.Set(keyValue[0], keyValue[1])
		}

		// Send request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			a.log.Errorf("Error sending request: %v", err)
			continue
		}
		defer resp.Body.Close()

		// Read response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			a.log.Errorf("Error reading response body: %v", err)
			continue
		}

		a.log.Debugf("Response status code: %d", resp.StatusCode)
		a.log.Debugf("Response body: %s", body)
	}
}

func main() {
	// Load API service from file
	file, err := os.Open("api_service.json")
	if err != nil {
		logrus.Fatal(err)
	}
	defer file.Close()

	var service APIService
	err = json.NewDecode(file, &service)
	if err != nil {
		logrus.Fatal(err)
	}

	// Create API analyzer
	analyzer := NewAPIAnalyzer()

	// Analyze API service
	analyzer.Analyze(&service)
}
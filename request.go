package main

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Data is used to store data from all the relevant endpoints in the API
type Data struct {
	Services []struct {
		// HealthState string `json:"health_check"`
		Name      string `json:"name"`
		State     string `json:"state"`
		Scale     int    `json:"scale"`
		StackName string `json:"stack_name"`
		System    bool   `json:"system"`
		EnvID     string `json:"environment_uuid"`
		Kind      string `json:"metadata_kind"`
	} `json:"services"`
	Stacks []struct {
		Name   string `json:"name"`
		Kind   string `json:"metadata_kind"`
		System bool   `json:"system"`
	}
	Hosts []struct {
		Name string `json:"name"`
		Kind string `json:"metadata_kind"`
	}
}

// gatherData - Collects the data from thw API, invokes functions to transform that data into metrics
func (e *Exporter) gatherData() (*Data, error) {

	// Create new data slice from Struct
	var data = new(Data)

	// Scrape EndPoint for JSON Data
	err := getJSON(e.MetaDataURL, &data)
	if err != nil {
		log.Errorf("Error getting JSON from URL: %s", err)
		return nil, err
	}
	log.Debugf("JSON Fetched : ", data)

	return data, err
}

// getJSON return json from server, return the formatted JSON
func getJSON(url string, target interface{}) error {

	log.Info("Scraping: ", url)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json")

	log.Info("header set")

	resp, err := client.Do(req)

	if err != nil {
		log.Error("Error Collecting JSON from API: ", err)
		panic(err)
	}

	respFormatted := json.NewDecoder(resp.Body).Decode(target)
	log.Info("formatted")

	// Close the response body, the underlying Transport should then close the connection.
	resp.Body.Close()

	// return formatted JSON
	return respFormatted
}

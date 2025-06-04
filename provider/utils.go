package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func sendPost(endpoint string, body map[string]interface{}, cookies map[string]string, bodyDest interface{}) (*http.Response, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	for name, value := range cookies {
		request.AddCookie(&http.Cookie{
			Name:  name,
			Value: value,
		})
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("did not get a 200 response")
	}

	json.NewDecoder(resp.Body).Decode(&bodyDest)

	//TODO: Check body for errors because we always get a 200

	return resp, nil
}

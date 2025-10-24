package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type HeihuRequest struct {
	ProductCode string `json:"productCode"`
}

// QueryHeihu queries the external Heihu API with the given item code
func QueryHeihu(itemCode string) (map[string]interface{}, error) {
	// Get configuration from environment variables
	heihuLink := os.Getenv("HEIHU_LINK")
	heihuSubLink := os.Getenv("HEIHU_SUB_LINK")
	xAuth := os.Getenv("X_AUTH")

	if heihuLink == "" || heihuSubLink == "" || xAuth == "" {
		return nil, fmt.Errorf("missing Heihu API configuration in environment variables")
	}

	// Build the full URL
	url := heihuLink + heihuSubLink

	// Create request body
	requestBody := HeihuRequest{
		ProductCode: itemCode,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth", xAuth)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Check if response status is not OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse response as generic JSON
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("error parsing response JSON: %v", err)
	}

	// Filter the results to only return exact matches
	if data, ok := result["data"].(map[string]interface{}); ok {
		if dataArray, ok := data["data"].([]interface{}); ok {
			// Find the exact match
			for _, item := range dataArray {
				if product, ok := item.(map[string]interface{}); ok {
					if productCode, ok := product["productCode"].(string); ok {
						if productCode == itemCode {
							// Return only the matched product
							return map[string]interface{}{
								"code": result["code"],
								"data": product,
								"msg":  result["msg"],
							}, nil
						}
					}
				}
			}
			// No exact match found
			return nil, fmt.Errorf("no product found with exact code: %s", itemCode)
		}
	}

	return result, nil
}
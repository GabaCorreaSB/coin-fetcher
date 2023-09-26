// This section specifies the package name and imports required for the code.
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"coinfetcher/types"
)

// Client represents a client for fetching cryptocurrency prices.
type Client struct {
	endpoint string // The endpoint URL for the price service.
}

// New creates a new instance of the Client with the specified endpoint.
func New(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
	}
}

// FetchPrice fetches cryptocurrency price information for the given ticker.
func (c *Client) FetchPrice(ctx context.Context, ticker string) (*types.PriceResponse, error) {
	// Create the full endpoint URL by combining the base endpoint and ticker parameter.
	endpoint := fmt.Sprintf("%s?ticker=%s", c.endpoint, ticker)

	// Create an HTTP GET request to the endpoint.
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	// Send the HTTP request using the default HTTP client.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Check if the response status code is not OK (200).
	if resp.StatusCode != http.StatusOK {
		// Decode the response body into a map for error information.
		httpErr := map[string]interface{}{}
		if err := json.NewDecoder(resp.Body).Decode(&httpErr); err != nil {
			return nil, err
		}
		// Return an error with the service response error message.
		return nil, fmt.Errorf("service responded with non-OK status code: %s", httpErr["error"])
	}

	// Create a new instance of PriceResponse to hold the decoded JSON response.
	priceResp := new(types.PriceResponse)
	if err := json.NewDecoder(resp.Body).Decode(priceResp); err != nil {
		return nil, err
	}

	// Return the successfully fetched PriceResponse.
	return priceResp, nil
}

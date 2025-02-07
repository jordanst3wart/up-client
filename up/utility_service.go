package up

import (
	"context"
	"fmt"
	"net/http"
)

// UtilityService handles communication with the utility related methods
type UtilityService service

// PingResponse represents the response from the ping endpoint
type PingResponse struct {
	Meta struct {
		ID          string `json:"id"`
		StatusEmoji string `json:"statusEmoji"`
	} `json:"meta"`
}

// Ping makes a basic ping request to verify authentication
func (s *UtilityService) Ping(ctx context.Context) (*PingResponse, *http.Response, error) {
	req, err := s.client.newRequest("GET", "util/ping", nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating request: %v", err)
	}

	var response PingResponse
	resp, err := s.client.do(ctx, req, &response)
	if err != nil {
		return nil, resp, err
	}

	return &response, resp, nil
}

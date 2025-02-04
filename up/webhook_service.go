package up

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// WebhooksService handles communication with the webhook related methods
type WebhooksService service

// WebhookEventTypeEnum represents the type of webhook event
type WebhookEventTypeEnum string

const (
	WebhookEventTransactionCreated WebhookEventTypeEnum = "TRANSACTION_CREATED"
	WebhookEventTransactionSettled WebhookEventTypeEnum = "TRANSACTION_SETTLED"
	WebhookEventTransactionDeleted WebhookEventTypeEnum = "TRANSACTION_DELETED"
	WebhookEventPing               WebhookEventTypeEnum = "PING"
)

// Webhook represents a webhook in Up
type Webhook struct {
	Type          string               `json:"type"`
	ID            string               `json:"id"`
	Attributes    WebhookAttributes    `json:"attributes"`
	Relationships WebhookRelationships `json:"relationships"`
	Links         Links                `json:"links"`
}

// WebhookAttributes represents the attributes of a webhook
type WebhookAttributes struct {
	URL         string    `json:"url"`
	Description *string   `json:"description"`
	SecretKey   string    `json:"secretKey,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

// WebhookRelationships represents the relationships of a webhook
type WebhookRelationships struct {
	Logs struct {
		Links Links `json:"links"`
	} `json:"logs"`
}

// WebhookListResponse represents the response from listing webhooks
type WebhookListResponse struct {
	Data  []Webhook `json:"data"`
	Links Links     `json:"links"`
}

// WebhookResponse represents the response for a single webhook
type WebhookResponse struct {
	Data Webhook `json:"data"`
}

// WebhookCreateRequest represents the request to create a webhook
type WebhookCreateRequest struct {
	Data struct {
		Attributes WebhookInputAttributes `json:"attributes"`
	} `json:"data"`
}

// WebhookInputAttributes represents the attributes for creating a webhook
type WebhookInputAttributes struct {
	URL         string  `json:"url"`
	Description *string `json:"description,omitempty"`
}

// WebhookDeliveryStatusEnum represents the status of a webhook delivery
type WebhookDeliveryStatusEnum string

const (
	WebhookDeliveryStatusDelivered       WebhookDeliveryStatusEnum = "DELIVERED"
	WebhookDeliveryStatusUndeliverable   WebhookDeliveryStatusEnum = "UNDELIVERABLE"
	WebhookDeliveryStatusBadResponseCode WebhookDeliveryStatusEnum = "BAD_RESPONSE_CODE"
)

// WebhookDeliveryLog represents a webhook delivery log
type WebhookDeliveryLog struct {
	Type          string                          `json:"type"`
	ID            string                          `json:"id"`
	Attributes    WebhookDeliveryLogAttributes    `json:"attributes"`
	Relationships WebhookDeliveryLogRelationships `json:"relationships"`
}

// WebhookDeliveryLogAttributes represents attributes of a webhook delivery log
type WebhookDeliveryLogAttributes struct {
	Request struct {
		Body string `json:"body"`
	} `json:"request"`
	Response *struct {
		StatusCode int    `json:"statusCode"`
		Body       string `json:"body"`
	} `json:"response"`
	DeliveryStatus WebhookDeliveryStatusEnum `json:"deliveryStatus"`
	CreatedAt      time.Time                 `json:"createdAt"`
}

// WebhookDeliveryLogRelationships represents relationships of a webhook delivery log
type WebhookDeliveryLogRelationships struct {
	WebhookEvent struct {
		Data struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		} `json:"data"`
	} `json:"webhookEvent"`
}

// WebhookDeliveryLogListResponse represents the response from listing webhook delivery logs
type WebhookDeliveryLogListResponse struct {
	Data  []WebhookDeliveryLog `json:"data"`
	Links Links                `json:"links"`
}

// WebhookEvent represents a webhook event
type WebhookEvent struct {
	Type          string                 `json:"type"`
	ID            string                 `json:"id"`
	Attributes    WebhookEventAttributes `json:"attributes"`
	Relationships struct {
		Webhook     WebhookRelationships      `json:"webhook"`
		Transaction *TransactionRelationships `json:"transaction,omitempty"`
	} `json:"relationships"`
}

// WebhookEventAttributes represents attributes of a webhook event
type WebhookEventAttributes struct {
	EventType WebhookEventTypeEnum `json:"eventType"`
	CreatedAt time.Time            `json:"createdAt"`
}

// WebhookEventResponse represents the response for a webhook event
type WebhookEventResponse struct {
	Data WebhookEvent `json:"data"`
}

// List returns a list of all webhooks
func (s *WebhooksService) List(ctx context.Context, opts *ListOptions) (*WebhookListResponse, *http.Response, error) {
	u := "webhooks"
	if opts != nil {
		var err error
		u, err = addOptions(u, opts)
		if err != nil {
			return nil, nil, err
		}
	}

	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var response WebhookListResponse
	resp, err := s.client.do(ctx, req, &response)
	if err != nil {
		return nil, resp, err
	}

	return &response, resp, nil
}

// Get returns a specific webhook by ID
func (s *WebhooksService) Get(ctx context.Context, webhookID string) (*Webhook, *http.Response, error) {
	u := fmt.Sprintf("webhooks/%s", webhookID)
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var response WebhookResponse
	resp, err := s.client.do(ctx, req, &response)
	if err != nil {
		return nil, resp, err
	}

	return &response.Data, resp, nil
}

// Create creates a new webhook
func (s *WebhooksService) Create(ctx context.Context, url string, description *string) (*Webhook, *http.Response, error) {
	u := "webhooks"

	createReq := &WebhookCreateRequest{}
	createReq.Data.Attributes = WebhookInputAttributes{
		URL:         url,
		Description: description,
	}

	req, err := s.client.newRequest("POST", u, createReq)
	if err != nil {
		return nil, nil, err
	}

	var response WebhookResponse
	resp, err := s.client.do(ctx, req, &response)
	if err != nil {
		return nil, resp, err
	}

	return &response.Data, resp, nil
}

// Delete deletes a webhook
func (s *WebhooksService) Delete(ctx context.Context, webhookID string) (*http.Response, error) {
	u := fmt.Sprintf("webhooks/%s", webhookID)
	req, err := s.client.newRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.do(ctx, req, nil)
}

// Ping sends a ping event to a webhook
func (s *WebhooksService) Ping(ctx context.Context, webhookID string) (*WebhookEvent, *http.Response, error) {
	u := fmt.Sprintf("webhooks/%s/ping", webhookID)
	req, err := s.client.newRequest("POST", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var response WebhookEventResponse
	resp, err := s.client.do(ctx, req, &response)
	if err != nil {
		return nil, resp, err
	}

	return &response.Data, resp, nil
}

// ListLogs returns a list of delivery logs for a webhook
func (s *WebhooksService) ListLogs(ctx context.Context, webhookID string, opts *ListOptions) (*WebhookDeliveryLogListResponse, *http.Response, error) {
	u := fmt.Sprintf("webhooks/%s/logs", webhookID)
	if opts != nil {
		var err error
		u, err = addOptions(u, opts)
		if err != nil {
			return nil, nil, err
		}
	}

	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var response WebhookDeliveryLogListResponse
	resp, err := s.client.do(ctx, req, &response)
	if err != nil {
		return nil, resp, err
	}

	return &response, resp, nil
}

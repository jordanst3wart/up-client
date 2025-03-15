package up

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	defaultBaseURL = "https://api.up.com.au/api/v1/"
	userAgent      = "up-go-client/1.0"
)

// Client manages communication with Up API
type Client struct {
	client  *http.Client
	baseURL *url.URL
	token   string

	common service // Reuse a single struct instead of creating one for each service

	// Services
	Accounts     *AccountsService
	Categories   *CategoriesService
	Tags         *TagsService
	Transactions *TransactionsService
	Webhooks     *WebhooksService
	Utility      *UtilityService
}

type service struct {
	client *Client
}

// NewClient returns a new Up API client
func NewClient(token string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: time.Second * 30,
		}
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		client:  httpClient,
		baseURL: baseURL,
		token:   token,
	}

	c.common.client = c

	// Initialize services
	c.Accounts = (*AccountsService)(&c.common)
	c.Categories = (*CategoriesService)(&c.common)
	c.Tags = (*TagsService)(&c.common)
	c.Transactions = (*TransactionsService)(&c.common)
	c.Webhooks = (*WebhooksService)(&c.common)
	c.Utility = (*UtilityService)(&c.common)

	return c
}

// newRequest creates an API request
func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	uri := c.baseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, uri.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	return req, nil
}

// do sends an API request and returns the API response
func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, err
	}

	defer resp.Body.Close()

	// probably a bug here somewhere
	if resp.StatusCode >= 400 {
		errorResponse := &ErrorResponse{}
		if err := json.NewDecoder(resp.Body).Decode(errorResponse); err != nil {
			return resp, fmt.Errorf("http status %d: failed to decode error response", resp.StatusCode)
		}
		return resp, errorResponse
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return resp, err
		}
	}

	return resp, err
}

// addOptions adds the parameters in opt as URL query parameters to s.
func addOptions(s string, opts interface{}) (string, error) {
	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opts)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// ListOptions specifies the optional parameters for pagination
type ListOptions struct {
	PageSize int    `url:"page[size],omitempty"`
	After    string `url:"page[after],omitempty"`
	Before   string `url:"page[before],omitempty"`
}

// Links represents pagination links
type Links struct {
	Prev string `json:"prev,omitempty"`
	Next string `json:"next,omitempty"`
}

// ErrorResponse represents an error response from the Up API
type ErrorResponse struct {
	Errors []ErrorObject `json:"errors"`
}

func (e *ErrorResponse) Error() string {
	if len(e.Errors) > 0 {
		return fmt.Sprintf("%s: %s", e.Errors[0].Title, e.Errors[0].Detail)
	}
	return "unknown error"
}

// ErrorObject represents a single error from the Up API
type ErrorObject struct {
	Status string `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Source *struct {
		Parameter string `json:"parameter,omitempty"`
		Pointer   string `json:"pointer,omitempty"`
	} `json:"source,omitempty"`
}

// MoneyObject represents a monetary value
type MoneyObject struct {
	CurrencyCode     string `json:"currencyCode"`
	Value            string `json:"value"`
	ValueInBaseUnits int64  `json:"valueInBaseUnits"`
}

// paginate handles pagination for any API resource that returns a list response with a Next link
func (c *Client) paginate(ctx context.Context, initialURL string, result interface{}) (*http.Response, error) {
	nextPageURL := initialURL

	// Get the value of the result pointer
	resultValue := reflect.ValueOf(result).Elem()

	// The first iteration is special since we need to initialize the result
	firstPage := true

	for nextPageURL != "" {
		req, err := c.newRequest("GET", nextPageURL, nil)
		if err != nil {
			return nil, err
		}

		// Create a new instance of the same type as result
		pageResult := reflect.New(resultValue.Type()).Interface()

		resp, err := c.do(ctx, req, pageResult)
		if err != nil {
			return resp, err
		}

		// Get the value of the page result
		pageValue := reflect.ValueOf(pageResult).Elem()

		if firstPage {
			// First page: just copy the entire result
			resultValue.Set(pageValue)
			firstPage = false
		} else {
			// Subsequent pages: append Data field contents
			resultData := resultValue.FieldByName("Data")
			pageData := pageValue.FieldByName("Data")

			if resultData.IsValid() && pageData.IsValid() {
				resultData.Set(reflect.AppendSlice(resultData, pageData))
			}
		}

		// Check for a Next link
		linksField := pageValue.FieldByName("Links")
		if !linksField.IsValid() {
			break
		}

		nextField := linksField.FieldByName("Next")
		if !nextField.IsValid() || nextField.String() == "" {
			break
		}

		nextPageURL = nextField.String()
	}

	return nil, nil
}

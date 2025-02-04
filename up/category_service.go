package up

import (
	"context"
	"fmt"
	"net/http"
)

// CategoriesService handles communication with the category related methods
type CategoriesService service

// Category represents a category in Up
type Category struct {
	Type          string                `json:"type"`
	ID            string                `json:"id"`
	Attributes    CategoryAttributes    `json:"attributes"`
	Relationships CategoryRelationships `json:"relationships"`
	Links         Links                 `json:"links"`
}

// CategoryAttributes represents the attributes of a category
type CategoryAttributes struct {
	Name string `json:"name"`
}

// CategoryRelationships represents the relationships of a category
type CategoryRelationships struct {
	Parent   CategoryParentRelationship   `json:"parent"`
	Children CategoryChildrenRelationship `json:"children"`
}

type CategoryParentRelationship struct {
	Data  *CategoryData `json:"data"`
	Links Links         `json:"links"`
}

type CategoryChildrenRelationship struct {
	Data  []CategoryData `json:"data"`
	Links Links          `json:"links"`
}

type CategoryData struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// CategoryListResponse represents the response from listing categories
type CategoryListResponse struct {
	Data []Category `json:"data"`
}

// CategoryGetResponse represents the response from getting a single category
type CategoryGetResponse struct {
	Data Category `json:"data"`
}

// ListCategoriesOptions specifies the optional parameters for listing categories
type ListCategoriesOptions struct {
	Parent string `url:"filter[parent],omitempty"`
}

// List returns a list of all categories
func (s *CategoriesService) List(ctx context.Context, opts *ListCategoriesOptions) (*CategoryListResponse, *http.Response, error) {
	u := "categories"
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

	var categories CategoryListResponse
	resp, err := s.client.do(ctx, req, &categories)
	if err != nil {
		return nil, resp, err
	}

	return &categories, resp, nil
}

// Get returns a specific category by ID
func (s *CategoriesService) Get(ctx context.Context, categoryID string) (*Category, *http.Response, error) {
	u := fmt.Sprintf("categories/%s", categoryID)
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var categoryResponse CategoryGetResponse
	resp, err := s.client.do(ctx, req, &categoryResponse)
	if err != nil {
		return nil, resp, err
	}

	return &categoryResponse.Data, resp, nil
}

// CategoryUpdateRequest represents the request to update a transaction's category
type CategoryUpdateRequest struct {
	Data *CategoryData `json:"data"`
}

// UpdateTransactionCategory updates the category of a transaction
func (s *CategoriesService) UpdateTransactionCategory(ctx context.Context, transactionID string, categoryID string) (*http.Response, error) {
	u := fmt.Sprintf("transactions/%s/relationships/category", transactionID)

	categoryData := &CategoryUpdateRequest{
		Data: &CategoryData{
			Type: "categories",
			ID:   categoryID,
		},
	}

	req, err := s.client.newRequest("PATCH", u, categoryData)
	if err != nil {
		return nil, err
	}

	return s.client.do(ctx, req, nil)
}

// RemoveTransactionCategory removes the category from a transaction
func (s *CategoriesService) RemoveTransactionCategory(ctx context.Context, transactionID string) (*http.Response, error) {
	u := fmt.Sprintf("transactions/%s/relationships/category", transactionID)

	categoryData := &CategoryUpdateRequest{
		Data: nil,
	}

	req, err := s.client.newRequest("PATCH", u, categoryData)
	if err != nil {
		return nil, err
	}

	return s.client.do(ctx, req, nil)
}

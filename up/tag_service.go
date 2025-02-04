package up

import (
	"context"
	"fmt"
	"net/http"
)

// TagsService handles communication with the tag related methods
type TagsService service

// Tag represents a tag in Up
type Tag struct {
	Type          string           `json:"type"`
	ID            string           `json:"id"`
	Relationships TagRelationships `json:"relationships"`
}

// TagRelationships represents the relationships of a tag
type TagRelationships struct {
	Transactions struct {
		Links Links `json:"links"`
	} `json:"transactions"`
}

// TagListResponse represents the response from listing tags
type TagListResponse struct {
	Data  []Tag `json:"data"`
	Links Links `json:"links"`
}

// TagInputResource represents a tag input for adding/removing tags
type TagInputResource struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// UpdateTransactionTagsRequest represents the request to update transaction tags
type UpdateTransactionTagsRequest struct {
	Data []TagInputResource `json:"data"`
}

// List returns a list of all tags
func (s *TagsService) List(ctx context.Context, opts *ListOptions) (*TagListResponse, *http.Response, error) {
	u := "tags"
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

	var tags TagListResponse
	resp, err := s.client.do(ctx, req, &tags)
	if err != nil {
		return nil, resp, err
	}

	return &tags, resp, nil
}

// AddToTransaction adds tags to a transaction
func (s *TagsService) AddToTransaction(ctx context.Context, transactionID string, tagIDs []string) (*http.Response, error) {
	u := fmt.Sprintf("transactions/%s/relationships/tags", transactionID)

	tags := make([]TagInputResource, len(tagIDs))
	for i, id := range tagIDs {
		tags[i] = TagInputResource{
			Type: "tags",
			ID:   id,
		}
	}

	req, err := s.client.newRequest("POST", u, &UpdateTransactionTagsRequest{Data: tags})
	if err != nil {
		return nil, err
	}

	return s.client.do(ctx, req, nil)
}

// RemoveFromTransaction removes tags from a transaction
func (s *TagsService) RemoveFromTransaction(ctx context.Context, transactionID string, tagIDs []string) (*http.Response, error) {
	u := fmt.Sprintf("transactions/%s/relationships/tags", transactionID)

	tags := make([]TagInputResource, len(tagIDs))
	for i, id := range tagIDs {
		tags[i] = TagInputResource{
			Type: "tags",
			ID:   id,
		}
	}

	req, err := s.client.newRequest("DELETE", u, &UpdateTransactionTagsRequest{Data: tags})
	if err != nil {
		return nil, err
	}

	return s.client.do(ctx, req, nil)
}

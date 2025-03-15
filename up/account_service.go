package up

import (
	"context"
	"fmt"
	"net/http"
)

// AccountsService handles communication with the account related methods
type AccountsService service

// AccountTypeEnum represents the type of bank account
type AccountTypeEnum string

const (
	AccountTypeSaver         AccountTypeEnum = "SAVER"
	AccountTypeTransactional AccountTypeEnum = "TRANSACTIONAL"
	AccountTypeHomeLoan      AccountTypeEnum = "HOME_LOAN"
)

// OwnershipTypeEnum represents the structure under which a bank account is owned
type OwnershipTypeEnum string

const (
	OwnershipTypeIndividual OwnershipTypeEnum = "INDIVIDUAL"
	OwnershipTypeJoint      OwnershipTypeEnum = "JOINT"
)

// Account represents an Up bank account
type Account struct {
	Type          string               `json:"type"`
	ID            string               `json:"id"`
	Attributes    AccountAttributes    `json:"attributes"`
	Relationships AccountRelationships `json:"relationships"`
	Links         Links                `json:"links"`
}

// AccountAttributes represents the attributes of an account
type AccountAttributes struct {
	DisplayName   string            `json:"displayName"`
	AccountType   AccountTypeEnum   `json:"accountType"`
	OwnershipType OwnershipTypeEnum `json:"ownershipType"`
	Balance       MoneyObject       `json:"balance"`
	CreatedAt     string            `json:"createdAt"`
}

// AccountRelationships represents the relationships of an account
type AccountRelationships struct {
	Transactions struct {
		Links Links `json:"links"`
	} `json:"transactions"`
}

// AccountListResponse represents the response from listing accounts
type AccountListResponse struct {
	Data  []Account `json:"data"`
	Links Links     `json:"links"`
}

// AccountGetResponse represents the response from getting a single account
type AccountGetResponse struct {
	Data Account `json:"data"`
}

// ListAccountsOptions specifies the optional parameters for listing accounts
type ListAccountsOptions struct {
	ListOptions
	AccountType   AccountTypeEnum   `url:"filter[accountType],omitempty"`
	OwnershipType OwnershipTypeEnum `url:"filter[ownershipType],omitempty"`
}

// List returns a list of all accounts
func (s *AccountsService) List(ctx context.Context, opts *ListAccountsOptions) (*AccountListResponse, *http.Response, error) {
	u := "accounts"
	if opts != nil {
		var err error
		u, err = addOptions(u, opts)
		if err != nil {
			return nil, nil, err
		}
	}

	var accounts AccountListResponse
	resp, err := s.client.paginate(ctx, u, &accounts)
	if err != nil {
		return nil, resp, err
	}

	return &accounts, resp, nil
}

// Get returns a specific account by ID
func (s *AccountsService) Get(ctx context.Context, accountID string) (*Account, *http.Response, error) {
	u := fmt.Sprintf("accounts/%s", accountID)
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var accountResponse AccountGetResponse
	resp, err := s.client.do(ctx, req, &accountResponse)
	if err != nil {
		return nil, resp, err
	}

	return &accountResponse.Data, resp, nil
}

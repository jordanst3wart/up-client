package up

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// TransactionsService handles communication with the transaction related methods
type TransactionsService service

// TransactionStatusEnum represents the status of a transaction
type TransactionStatusEnum string

const (
	TransactionStatusHeld    TransactionStatusEnum = "HELD"
	TransactionStatusSettled TransactionStatusEnum = "SETTLED"
)

// Transaction represents a transaction in Up
type Transaction struct {
	Type          string                   `json:"type"`
	ID            string                   `json:"id"`
	Attributes    TransactionAttributes    `json:"attributes"`
	Relationships TransactionRelationships `json:"relationships"`
	Links         Links                    `json:"links"`
}

// TransactionAttributes represents the attributes of a transaction
type TransactionAttributes struct {
	Status             TransactionStatusEnum `json:"status"`
	RawText            *string               `json:"rawText"`
	Description        string                `json:"description"`
	Message            *string               `json:"message"`
	IsCategorizable    bool                  `json:"isCategorizable"`
	HoldInfo           *HoldInfo             `json:"holdInfo"`
	RoundUp            *RoundUp              `json:"roundUp"`
	Cashback           *Cashback             `json:"cashback"`
	Amount             MoneyObject           `json:"amount"`
	ForeignAmount      *MoneyObject          `json:"foreignAmount"`
	CardPurchaseMethod *CardPurchaseMethod   `json:"cardPurchaseMethod"`
	SettledAt          *time.Time            `json:"settledAt"`
	CreatedAt          time.Time             `json:"createdAt"`
}

type HoldInfo struct {
	Amount        MoneyObject  `json:"amount"`
	ForeignAmount *MoneyObject `json:"foreignAmount"`
}

type RoundUp struct {
	Amount       MoneyObject  `json:"amount"`
	BoostPortion *MoneyObject `json:"boostPortion"`
}

type Cashback struct {
	Description string      `json:"description"`
	Amount      MoneyObject `json:"amount"`
}

// CardPurchaseMethodEnum represents the type of card purchase
type CardPurchaseMethodEnum string

const (
	CardPurchaseBarCode        CardPurchaseMethodEnum = "BAR_CODE"
	CardPurchaseOCR            CardPurchaseMethodEnum = "OCR"
	CardPurchaseCardPin        CardPurchaseMethodEnum = "CARD_PIN"
	CardPurchaseCardDetails    CardPurchaseMethodEnum = "CARD_DETAILS"
	CardPurchaseCardOnFile     CardPurchaseMethodEnum = "CARD_ON_FILE"
	CardPurchaseEcommerce      CardPurchaseMethodEnum = "ECOMMERCE"
	CardPurchaseMagneticStripe CardPurchaseMethodEnum = "MAGNETIC_STRIPE"
	CardPurchaseContactless    CardPurchaseMethodEnum = "CONTACTLESS"
)

type CardPurchaseMethod struct {
	Method           CardPurchaseMethodEnum `json:"method"`
	CardNumberSuffix *string                `json:"cardNumberSuffix"`
}

// TransactionRelationships represents the relationships of a transaction
type TransactionRelationships struct {
	Account         AccountRelationship   `json:"account"`
	TransferAccount *AccountRelationship  `json:"transferAccount"`
	Category        *CategoryRelationship `json:"category"`
	ParentCategory  *CategoryRelationship `json:"parentCategory"`
	Tags            TagRelationships      `json:"tags"`
}

type AccountRelationship struct {
	Data  AccountData `json:"data"`
	Links Links       `json:"links,omitempty"`
}

type AccountData struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type CategoryRelationship struct {
	Data  *CategoryData `json:"data"`
	Links Links         `json:"links,omitempty"`
}

// TransactionListResponse represents the response from listing transactions
type TransactionListResponse struct {
	Data  []Transaction `json:"data"`
	Links Links         `json:"links"`
}

// TransactionGetResponse represents the response from getting a single transaction
type TransactionGetResponse struct {
	Data Transaction `json:"data"`
}

// ListTransactionsOptions specifies the optional parameters for listing transactions
type ListTransactionsOptions struct {
	ListOptions
	Status   TransactionStatusEnum `url:"filter[status],omitempty"`
	Since    *time.Time            `url:"filter[since],omitempty"`
	Until    *time.Time            `url:"filter[until],omitempty"`
	Category string                `url:"filter[category],omitempty"`
	Tag      string                `url:"filter[tag],omitempty"`
}

// List returns a list of all transactions
func (s *TransactionsService) List(ctx context.Context, opts *ListTransactionsOptions) (*TransactionListResponse, *http.Response, error) {
	u := "transactions"
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

	var transactions TransactionListResponse
	resp, err := s.client.do(ctx, req, &transactions)
	if err != nil {
		return nil, resp, err
	}

	return &transactions, resp, nil
}

// ListByAccount returns a list of transactions for a specific account
func (s *TransactionsService) ListByAccount(ctx context.Context, accountID string, opts *ListTransactionsOptions) (*TransactionListResponse, *http.Response, error) {
	u := fmt.Sprintf("accounts/%s/transactions", accountID)
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

	var transactions TransactionListResponse
	resp, err := s.client.do(ctx, req, &transactions)
	if err != nil {
		return nil, resp, err
	}

	return &transactions, resp, nil
}

// Get returns a specific transaction by ID
func (s *TransactionsService) Get(ctx context.Context, transactionID string) (*Transaction, *http.Response, error) {
	u := fmt.Sprintf("transactions/%s", transactionID)
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var transactionResponse TransactionGetResponse
	resp, err := s.client.do(ctx, req, &transactionResponse)
	if err != nil {
		return nil, resp, err
	}

	return &transactionResponse.Data, resp, nil
}

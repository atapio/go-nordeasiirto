package nordeasiirto

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	paymentBasePath = "/payment"
)

type PaymentService interface {
	SendIBANPayment(context.Context, *PaymentRequest) (*PaymentResponse, *Response, error)
	GetStatus(context.Context, string) (*PaymentStatusResponse, *Response, error)
}

type PaymentServiceOp struct {
	client *Client
}

var _ PaymentService = &PaymentServiceOp{}

type PaymentRequest struct {
	LookupID              string   `json:"lookupId"`
	Amount                int      `json:"amount"`
	Currency              string   `json:"currency"`
	BeneAccountNumber     string   `json:"beneAccountNumber"`
	FallbackPayment       string   `json:"fallbackPayment,omitempty"`
	BeneFirstNames        []string `json:"beneFirstNames,omitempty"`
	BeneLastName          string   `json:"beneLastName,omitempty"`
	BeneCompanyName       string   `json:"beneCompanyName,omitempty"`
	ReferenceNumber       string   `json:"referenceNumber,omitempty"`
	PaymentMessage        string   `json:"paymentMessage,omitempty"`
	BenecifiaryMinimumAge int      `json:"beneficiaryMinimumAge,omitempty"`
	BenecifiaryIdentifier string   `json:"beneficiaryIdentifier,omitempty"`
}

type PaymentResponse struct {
	Status           string    `json:"status"`
	ArchiveReference string    `json:"archiveReference"`
	FallbackPayment  bool      `json:"fallbackPayment"`
	Timestamp        time.Time `json:"timestamp"`
}

type PaymentStatusResponse struct {
	Status           string    `json:"status"`
	ArchiveReference string    `json:"archiveReference"`
	FallbackPayment  bool      `json:"fallbackPaymet"`
	PaymentTime      time.Time `json:"paymentTime"`
}

func (s *PaymentServiceOp) SendIBANPayment(ctx context.Context, pr *PaymentRequest) (*PaymentResponse, *Response, error) {
	if pr == nil {
		return nil, nil, NewArgError("paymentRequest", "cannot be nil")
	}

	path := paymentBasePath + "/pay"

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, pr)
	if err != nil {
		return nil, nil, err
	}

	paymentResponse := new(PaymentResponse)
	resp, err := s.client.Do(ctx, req, paymentResponse)
	if err != nil {
		return nil, resp, err
	}

	return paymentResponse, resp, nil
}

func (s *PaymentServiceOp) GetStatus(ctx context.Context, lookupID string) (*PaymentStatusResponse, *Response, error) {
	if lookupID == "" {
		return nil, nil, NewArgError("lookupID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/payment-status/%s", paymentBasePath, lookupID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	paymentStatusResponse := new(PaymentStatusResponse)
	resp, err := s.client.Do(ctx, req, paymentStatusResponse)
	if err != nil {
		return nil, resp, err
	}

	return paymentStatusResponse, resp, nil

}

package nordeasiirto

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	lookupBasePath = "/lookup"
)

type LookupService interface {
	GetUUID(context.Context) (*Lookup, *Response, error)
}

type Lookup struct {
	ID      string    `json:"lookupId"`
	Expires time.Time `json:"expires"`
}

type LookupServiceOp struct {
	client *Client
}

// assert type-correctness
var _ LookupService = &LookupServiceOp{}

// GETUUID fetches a Lookup ID for the payment
func (s *LookupServiceOp) GetUUID(ctx context.Context) (*Lookup, *Response, error) {
	path := fmt.Sprintf("%s/uuid", lookupBasePath)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	lookup := new(Lookup)
	resp, err := s.client.Do(ctx, req, lookup)
	if err != nil {
		return nil, resp, err
	}
	return lookup, resp, nil
}

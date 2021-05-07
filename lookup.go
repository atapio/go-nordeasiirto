package nordeasiirto

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	LookupBasePath = "/lookup"
)

type LookupService interface {
	Get(context.Context) (*Lookup, *Response, error)
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

// Get a Lookup UUID.
func (s *LookupServiceOp) Get(ctx context.Context) (*Lookup, *Response, error) {
	path := fmt.Sprintf("%s/uuid", LookupBasePath)

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

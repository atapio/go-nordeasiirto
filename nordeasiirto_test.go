package nordeasiirto

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	mux *http.ServeMux

	ctx = context.TODO()

	client *Client

	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewClient(nil)
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Errorf("Request method = %v, expected %v", r.Method, expected)
	}
}

func TestCheckResponse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		title    string
		input    *http.Response
		expected *ErrorResponse
	}{
		{
			title: "default",
			input: &http.Response{
				Request:    &http.Request{},
				StatusCode: http.StatusBadRequest,
				Body:       ioutil.NopCloser(strings.NewReader(`{"message":"m"}`)),
			},
			expected: &ErrorResponse{
				Message: "m",
			},
		},
		// ensure that we properly handle API errors that do not contain a
		// response body
		{
			title: "no body",
			input: &http.Response{
				Request:    &http.Request{},
				StatusCode: http.StatusBadRequest,
				Body:       ioutil.NopCloser(strings.NewReader("")),
			},
			expected: &ErrorResponse{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			err := CheckResponse(tc.input).(*ErrorResponse)
			require.Error(t, err)

			tc.expected.Response = tc.input
			assert.Equal(t, tc.expected, err)
		})
	}
}

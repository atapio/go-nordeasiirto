package nordeasiirto

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLookupGetUUID(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/lookup/uuid", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		response := `
		{
			"lookupId": "f6887f21-b824-4438-9ef6-d91e53d922fb",
			"expires": "2021-05-08T07:05:32Z"
			}
		}`

		fmt.Fprint(w, response)
	})

	lookup, _, err := client.Lookup.GetUUID(ctx)
	require.NoError(t, err, "Lookup.GetUUID returned error")

	expected := &Lookup{
		ID:      "f6887f21-b824-4438-9ef6-d91e53d922fb",
		Expires: time.Date(2021, 5, 8, 7, 5, 32, 0, time.UTC),
	}

	assert.Equal(t, expected, lookup)

}

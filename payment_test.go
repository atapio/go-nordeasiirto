package nordeasiirto

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetStatus(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/payment/payment-status/487c7a7e-bd1d-4791-8431-83d0755b8874", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		response := `
		{
			"status": "PAID",
			"archiveReference": "archiveref-123",
			"fallbackPayment": false,
			"paymentTime": "2021-05-08T07:05:32Z"
			}
		}`

		fmt.Fprint(w, response)
	})

	status, _, err := client.Payment.GetStatus(ctx, "487c7a7e-bd1d-4791-8431-83d0755b8874")
	require.NoError(t, err, "Payment.GetStatus returned error")

	expected := &PaymentStatusResponse{
		Status:           "PAID",
		ArchiveReference: "archiveref-123",
		FallbackPayment:  false,
		PaymentTime:      time.Date(2021, 5, 8, 7, 5, 32, 0, time.UTC),
	}

	assert.Equal(t, expected, status)
}

func TestSendIBANPayment(t *testing.T) {
	tests := []struct {
		title            string
		paymentRequest   *PaymentRequest
		response         string
		expectedRequest  map[string]interface{}
		expectedResponse *PaymentStatusResponse
	}{
		{
			title: "company",
			paymentRequest: &PaymentRequest{
				LookupID:          "01db84ee-15f8-4b82-9f77-7f6007b2ad77",
				Amount:            1000,
				Currency:          "EUR",
				BeneAccountNumber: "FI3815723500045661",
				BeneCompanyName:   "Integration company Oy",
				ReferenceNumber:   "RF111232",
				PaymentMessage:    "IBAN payment",
			},
			response: `
				{
					"status": "PAID",
					"archiveReference": "archiveref-123",
					"fallbackPayment": false,
					"paymentTime": "2021-05-08T07:05:32Z"
				}`,
			expectedRequest: map[string]interface{}{
				"lookupId":          "01db84ee-15f8-4b82-9f77-7f6007b2ad77",
				"amount":            float64(1000),
				"currency":          "EUR",
				"beneCompanyName":   "Integration company Oy",
				"beneAccountNumber": "FI3815723500045661",
				"referenceNumber":   "RF111232",
				"paymentMessage":    "IBAN payment",
			},
			expectedResponse: &PaymentStatusResponse{
				Status:           "PAID",
				ArchiveReference: "archiveref-123",
				FallbackPayment:  false,
				PaymentTime:      time.Date(2021, 5, 8, 7, 5, 32, 0, time.UTC),
			},
		},
		{
			title: "person",
			paymentRequest: &PaymentRequest{
				LookupID:          "d289afdf-29e6-4452-86d8-e61e14213d25",
				Amount:            1000,
				Currency:          "EUR",
				BeneAccountNumber: "FI3815723500045661",
				BeneLastName:      "Kankkunen",
				BeneFirstNames:    []string{"Timo", "Pekka"},
				ReferenceNumber:   "RF111232",
				PaymentMessage:    "IBAN payment",
			},
			response: `
				{
					"status": "PROCESSING",
					"archiveReference": "archiveref-123",
					"fallbackPayment": false
				}`,
			expectedRequest: map[string]interface{}{
				"lookupId":          "d289afdf-29e6-4452-86d8-e61e14213d25",
				"amount":            float64(1000),
				"currency":          "EUR",
				"beneAccountNumber": "FI3815723500045661",
				"beneLastName":      "Kankkunen",
				"beneFirstNames":    []interface{}{"Timo", "Pekka"},
				"referenceNumber":   "RF111232",
				"paymentMessage":    "IBAN payment",
			},
			expectedResponse: &PaymentStatusResponse{
				Status:           "PROCESSING",
				ArchiveReference: "archiveref-123",
				FallbackPayment:  false,
				PaymentTime:      time.Time{},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			setup()
			defer teardown()

			mux.HandleFunc("/payment/pay", func(w http.ResponseWriter, r *http.Request) {
				var v map[string]interface{}
				err := json.NewDecoder(r.Body).Decode(&v)
				require.NoError(t, err)

				assert.Equal(t, tc.expectedRequest, v)

				fmt.Fprint(w, tc.response)
			})

			status, _, err := client.Payment.SendIBANPayment(ctx, tc.paymentRequest)
			require.NoError(t, err, "Payment.SendIBANPayment returned error")

			assert.Equal(t, tc.expectedResponse, status)

		})
	}
}

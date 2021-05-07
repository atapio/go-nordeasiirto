package nordeasiirto_test

import (
	"context"
	"testing"
	"time"

	"github.com/atapio/nordeasiirto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSuccessfulPayment(t *testing.T) {
	ctx := context.Background()
	config := &nordeasiirto.Config{
		Username:    "ME35912345671",
		Password:    "ee582575-7941-4d5a-b9b9-7a5101c2e2bc",
		Environment: nordeasiirto.Test,
	}

	client, err := nordeasiirto.NewFromConfig(ctx, config)
	assert.NoError(t, err)

	lookup, _, err := client.Lookup.Get(ctx)
	assert.NoError(t, err)

	req := &nordeasiirto.PaymentRequest{
		LookupID:          lookup.ID,
		BeneAccountNumber: "FI3815723500045661",
		BeneCompanyName:   "Acme Inc.",
		Amount:            1000,
		Currency:          "EUR",
	}

	paymentResp, _, err := client.Payment.SendIBANPayment(ctx, req)
	require.NoError(t, err)

	assert.Equal(t, "PAID", paymentResp.Status)
	assert.NotEmpty(t, paymentResp.ArchiveReference)
	assert.False(t, paymentResp.FallbackPayment)
	assert.NotEqual(t, paymentResp.Timestamp, time.Time{})

	statusResp, _, err := client.Payment.GetStatus(ctx, lookup.ID)
	require.NoError(t, err)

	assert.Equal(t, "PAID", statusResp.Status)
	assert.NotEmpty(t, statusResp.ArchiveReference)
	assert.False(t, statusResp.FallbackPayment)
	assert.NotEqual(t, statusResp.PaymentTime, time.Time{})
}

func TestUnsuccessfulPayment(t *testing.T) {
	ctx := context.Background()
	config := &nordeasiirto.Config{
		Username:    "ME35912345671",
		Password:    "ee582575-7941-4d5a-b9b9-7a5101c2e2bc",
		Environment: nordeasiirto.Test,
	}

	client, err := nordeasiirto.NewFromConfig(ctx, config)
	assert.NoError(t, err)

	lookup, _, err := client.Lookup.Get(ctx)
	assert.NoError(t, err)

	req := &nordeasiirto.PaymentRequest{
		LookupID:          lookup.ID,
		BeneAccountNumber: "FI3815723500045661",
		Amount:            1000,
		Currency:          "EUR",
	}

	_, _, err = client.Payment.SendIBANPayment(ctx, req)
	assert.Error(t, err, "expected 400")
	assert.Contains(t, err.Error(), "Bad Request - beneficiary name not valid or missing")

	_, _, err = client.Payment.GetStatus(ctx, lookup.ID)
	assert.Error(t, err, "expected 404")
	assert.Contains(t, err.Error(), "Not Found")
}

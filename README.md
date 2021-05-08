# Nordeasiirto

Nordeasiirto is a Go client library for accessing the [Nordea Siirto for Corporates API](https://www.nordea.fi/Images/146-261077/Nordea_Siirto_for_Corporates_API_specification_(06.2020).pdf).

**Note:** Currently only the IBAN payments are implemented.

## Usage

```go
ctx := context.Background()
config := &nordeasiirto.Config{
    Username:    "ME35912345671",
    Password:    "ee582575-7941-4d5a-b9b9-7a5101c2e2bc",
    Environment: nordeasiirto.Test,
}

client, err := nordeasiirto.NewFromConfig(ctx, config)
if err != nil {
   log.Fatal(err)
}

lookup, _, err := client.Lookup.GetUUID(ctx)
if err != nil {
   log.Fatal(err)
}

req := &nordeasiirto.PaymentRequest{
    LookupID:          lookup.ID,
    BeneAccountNumber: "FI3815723500045661",
    BeneCompanyName:   "Acme Inc.",
    Amount:            1000,
    Currency:          "EUR",
}

paymentResp, _, err := client.Payment.SendIBANPayment(ctx, req)
if err != nil {
   log.Fatal(err)
}

statusResp, _, err := client.Payment.GetStatus(ctx, lookup.ID)
if err != nil {
   log.Fatal(err)
}
```

## Acknowledgements

This library is based on the excellent [DigitalOcean Godo](https://github.com/digitalocean/godo) library.

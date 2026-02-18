# Payment Gateway Client (Pay / Verify)

### A small Go package that lets your service communicate with a payment gateway using two operations:

- **Pay**: create a payment request and receive a `referenceID` + `paymentURL`
- **Verify**: verify a payment after the user returns from the gateway

## Features
- Simple interface-based design (`Payment`)
- Context-aware calls (`context.Context`)
- Minimal request structures for Pay and Verify

---

## Installation

```bash
go get github.com/amirzayi/payment-go
```

Core Interface
```go
type Payment interface {
	Pay(ctx context.Context, in PayRequest) (referenceID string, paymentURL string, err error)
	Verify(ctx context.Context, in VerifiyRequest) error
}
```

PayRequest
```go
type PayRequest struct {
	Amount  uint64
	Email   string
	Mobile  string
	OrderID string
	UserID  string
}
```

Fields
* Amount (uint64): Amount to charge (use the gateway’s expected unit)
* Email (string): Customer email (optional depending on gateway)
* Mobile (string): Customer mobile (optional depending on gateway)
* OrderID (string): Your unique order identifier
* UserID (string): Your internal user/customer identifier

VerifiyRequest
```go
type VerifiyRequest struct {
	OrderID     string
	ReferenceID string
	Amount      uint64
}
```

Fields
* OrderID (string): Your order identifier
* ReferenceID (string): Gateway reference returned from Pay
* Amount (uint64): Amount you expect to verify (should match Pay)


Usage
The httpClient accept `*http.Client` that could be `http.DefaultClient` or custom client like sentry http client.

```go
package main

import (
	"context"
	"fmt"
	"log"

    paymentgo "github.com/amirzayi/payment-go"
	"github.com/amirzayi/payment-go/behpardakht"
    "github.com/amirzayi/payment-go/novinpay"
)

func main() {
	// Initialize gateway client either novinpay or behpardakht
    var pgw paymentgo.Payment

    npgwUserName := os.Getenv("NOVINPAY_USERNAME")
    npgwPassword := os.Getenv("NOVINPAY_PASSOWRD")
    npgwMerchantID := os.Getenv("NOVINPAY_MERCHANTID")
    npgwTerminalID := os.Getenv("NOVINPAY_TERMINALID")
    npmRedirectURL := os.Getenv("NOVINPAY_REDIRECT_URL")
    npgwCertificatePassword := os.Getenv("NOVINPAY_CERTIFICATE_PASSOWRD")
    f, err := os.Open("cert.p12")
    if err != nil {
        log.Fatalf("failed to load certificate file: %v", err)
    }
    pgw, err := novinpay.NewService(http.DefaultClient, novinpay.ServiceURL, novinpay.PaymentGatewayURL, npgwUserName, npgwPassword, npgwMerchantID, npgwTerminalID, npmRedirectURL, npgwCertificatePassword, f)
    if err != nil {
        log.Fatal(err)
    }

    bpgwUserName := os.Getenv("BEHPARDAKHT_USERNAME")
    bpgwPassword := os.Getenv("BEHPARDAKHT_PASSOWRD")
    bpgwTerminalID, _ := strconv.ParseInt(os.Getenv("BEHPARDAKHT_TERMINALID"), 10, 32)
    bpmRedirectURL := os.Getenv("BEHPARDAKHT_REDIRECT_URL")
    pgw = behpardakht.NewService(http.DefaultClient, behpardakht.ServiceURL, behpardakht.GatewayURL, bpgwUserName, bpgwPassword, bpmRedirectURL, int(bpgwTerminalID))

    orderID := "10001"
    amount := 150000
	refID, payURL, err := pgw.Pay(context.Background(), payment.PayRequest{
		Amount:  amount,
		Email:   "user@example.com",
		Mobile:  "09123456789",
		OrderID: orderID,
		UserID:  "42",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Redirect user to payURL in your HTTP handler.

    err = pgw.Verify(ctx, paymentgo.VerifiyRequest{
        OrderID:     orderID,
        ReferenceID: refID,
        Amount:      amount,
    })
}
```


Suggested Payment Flow (High Level)
* Client requests payment (your API)
* Your server calls Pay → gets (referenceID, paymentURL)
* Redirect user to paymentURL
* Gateway redirects user back to your callback
* Your server calls Verify with OrderID, referenceID, and expected Amount
* Update order status based on Verify result
package behpardakht

import (
	"context"
	"net/http"
	"strconv"

	paymentgo "github.com/amirzayi/payment-go"
)

const (
	ServiceURL = "https://bpm.shaparak.ir/pgwchannel/services/pgw?wsdl"
	GatewayURL = "https://bpm.shaparak.ir/pgwchannel/startpay.mellat"
)

type credentials struct {
	TerminalID   int    `xml:"terminalId"`   // شماره پايانه پذيرنده
	UserName     string `xml:"userName"`     // نام کاربری پذيرنده
	UserPassword string `xml:"userPassword"` // کلمه عبور پذيرنده
}

type service struct {
	httpClient  *http.Client
	serviceURL  string
	gatewayURL  string
	username    string
	password    string
	callbackURL string
	terminalID  int
}

func NewService(httpClient *http.Client, serviceURL, gatewayURL, username, password, callbackURL string, terminalID int) paymentgo.Payment {
	return service{
		httpClient:  httpClient,
		serviceURL:  serviceURL,
		gatewayURL:  gatewayURL,
		username:    username,
		password:    password,
		callbackURL: callbackURL,
		terminalID:  terminalID,
	}
}

func (s service) Verify(ctx context.Context, in paymentgo.VerifiyRequest) error {
	orderID, err := strconv.ParseUint(in.OrderID, 10, 64)
	if err != nil {
		return err
	}
	saleReferenceID, err := strconv.ParseUint(in.ReferenceID, 10, 64)
	if err != nil {
		return err
	}
	err = s.verify(ctx, verifyRequestData{
		OrderID:         orderID,
		SaleOrderID:     orderID,
		SaleReferenceID: saleReferenceID,
	})
	if err != nil {
		return err
	}
	err = s.Settle(ctx, settleRequestData{
		OrderID:         orderID,
		SaleOrderID:     orderID,
		SaleReferenceID: saleReferenceID,
	})
	return err
}

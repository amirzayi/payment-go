package novinpay

import (
	"context"
	"strings"

	paymentgo "github.com/amirzayi/payment-go"
	"github.com/amirzayi/payment-go/novinpay/base"
)

type verifyTransactionRequest struct {
	WsContext wsContext `json:"WSContext"`
	Token     string    `json:"Token"`
	RefNum    string    `json:"RefNum"`
}

type verifyResponse struct {
	Result string `json:"Result"`
	RefNum string `json:"RefNum"`
	Amount uint64 `json:"Amount"`
}

func (s service) Verify(ctx context.Context, in paymentgo.VerifiyRequest) error {
	txResult, err := base.DoPostApiCall[verifyResponse](
		ctx,
		s.serviceURL+verifyTransactionURL,
		verifyTransactionRequest{
			WsContext: wsContext{
				UserId:   s.userName,
				Password: s.password,
			},
			Token:  in.ReferenceID,
			RefNum: in.OrderID,
		},
	)
	if err != nil {
		return err
	}
	if !strings.EqualFold(txResult.Result, base.ResponseSuccess) {
		return base.GetResponseError(txResult.Result)
	}
	if in.OrderID != txResult.RefNum {
		return base.ErrMismatchVerificationRefnum
	}
	if in.Amount != txResult.Amount {
		return base.ErrMismatchVerificationAmount
	}
	return nil
}

type PaymentCallback struct {
	State             string `json:"State"`
	ResNum            string `json:"ResNum"` // Order ID
	RefNum            string `json:"RefNum"`
	CustomerRefNum    string `json:"CustomerRefNum"`
	MID               string `json:"MID"`
	Language          string `json:"Language"`
	CardHashPan       string `json:"CardHashPan"`
	CardMaskPan       string `json:"CardMaskPan"`
	GoodReferenceId   string `json:"GoodReferenceId"`
	MerchantData      string `json:"MerchantData"`
	TraceNo           string `json:"TraceNo"`
	Token             string `json:"token"`
	TransactionAmount string `json:"transactionAmount"`
	Email             string `json:"emailAddress"`
	Mobile            string `json:"mobileNo"`
}

package novinpay

import (
	"context"

	paymentgo "github.com/amirzayi/payment-go"
)

type payRequest struct {
	amount  uint64
	email   string
	mobile  string
	orderID string
}

func transformPayRequest(in paymentgo.PayRequest) payRequest {
	return payRequest{
		orderID: in.OrderID,
		amount:  in.Amount,
		mobile:  in.Mobile,
		email:   in.Email,
	}
}

func (s service) Pay(ctx context.Context, in paymentgo.PayRequest) (string, string, error) {
	_, err := s.login(ctx)
	if err != nil {
		return "", "", err
	}
	tx, err := s.generateTransactionDataToSign(ctx, transformPayRequest(in))
	if err != nil {
		return "", "", err
	}
	signature, err := s.signToken(tx.DataToSign)
	if err != nil {
		return "", "", err
	}
	response, err := s.generateSignedDataToken(ctx, signature, tx.UniqueId)
	if err != nil {
		return "", "", err
	}
	return response.Token, s.paymentGatewayURL, nil
}

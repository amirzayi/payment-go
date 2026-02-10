package paymentgo

import "context"

type Payment interface {
	Pay(ctx context.Context, in PayRequest) (referenceID string, paymentURL string, err error)
	Verify(ctx context.Context, in VerifiyRequest) error
}

type PayRequest struct {
	Amount  uint64
	Email   string
	Mobile  string
	OrderID string
	UserID  string
}

type VerifiyRequest struct {
	OrderID     string
	ReferenceID string
	Amount      uint64
}

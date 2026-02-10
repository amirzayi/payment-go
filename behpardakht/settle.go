package behpardakht

import (
	"context"
	"encoding/xml"

	"github.com/amirzayi/payment-go/behpardakht/base"
)

type settleRequestBody struct {
	XMLName         xml.Name      `xml:"soapenv:Body"`
	BPSettleRequest settleRequest `xml:"web:bpSettleRequest"`
}

type settleRequest struct {
	XMLName xml.Name `xml:"web:bpSettleRequest"`
	credentials
	settleRequestData
}

type settleRequestData struct {
	OrderID         uint64 `xml:"orderId"`         // شماره درخواست(واريز)
	SaleOrderID     uint64 `xml:"saleOrderId"`     // شماره درخواست خريد
	SaleReferenceID uint64 `xml:"saleReferenceId"` // کد مرجع تراکنش خريد
}

type settleResponse struct {
	Result base.Response `xml:"bpSettleRequestResponse"`
}

func (m service) Settle(ctx context.Context, in settleRequestData) error {
	response, err := base.Call[settleResponse](ctx, m.serviceURL, settleRequestBody{
		BPSettleRequest: settleRequest{
			credentials: credentials{
				TerminalID:   m.terminalID,
				UserName:     m.username,
				UserPassword: m.password,
			},
			settleRequestData: settleRequestData{
				OrderID:         in.OrderID,
				SaleOrderID:     in.SaleOrderID,
				SaleReferenceID: in.SaleReferenceID,
			},
		},
	})
	if err != nil {
		return err
	}
	return base.ConvertError(response.Result.Return)
}

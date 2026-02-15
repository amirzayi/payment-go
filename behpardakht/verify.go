package behpardakht

import (
	"context"
	"encoding/xml"

	"github.com/amirzayi/payment-go/behpardakht/base"
)

type verifyRequestBody struct {
	XMLName         xml.Name      `xml:"soapenv:Body"`
	BPVerifyRequest verifyRequest `xml:"web:bpVerifyRequest"`
}

type verifyRequest struct {
	XMLName xml.Name `xml:"web:bpVerifyRequest"`
	credentials
	verifyRequestData
}

type verifyRequestData struct {
	OrderID         uint64 `xml:"orderId"`         // شماره درخواست(استعلام)
	SaleOrderID     uint64 `xml:"saleOrderId"`     // شماره درخواست خريد
	SaleReferenceID uint64 `xml:"saleReferenceId"` // کد مرجع تراکنش خريد
}

type verifyResponse struct {
	Result base.Response `xml:"bpVerifyRequestResponse"`
}

func (m service) verify(ctx context.Context, in verifyRequestData) error {
	response, err := base.DoPostApiCall[verifyResponse](ctx, m.serviceURL, verifyRequestBody{
		BPVerifyRequest: verifyRequest{
			credentials: credentials{
				TerminalID:   m.terminalID,
				UserName:     m.username,
				UserPassword: m.password,
			},
			verifyRequestData: verifyRequestData{
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

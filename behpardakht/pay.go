package behpardakht

import (
	"context"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"

	paymentgo "github.com/amirzayi/payment-go"
	"github.com/amirzayi/payment-go/behpardakht/base"
)

type payResponse struct {
	Result base.Response `xml:"bpPayRequestResponse"`
}

type payRequestBody struct {
	XMLName      xml.Name   `xml:"soapenv:Body"`
	BpPayRequest payRequest `xml:"web:bpPayRequest"`
}

type payRequest struct {
	XMLName xml.Name `xml:"web:bpPayRequest"`
	credentials
	LocalDate   string `xml:"localDate"`   // تاريخ درخواست YYYYMMDD
	LocalTime   string `xml:"localTime"`   // ساعت درخواست HHMMSS
	CallbackURL string `xml:"callBackUrl"` // آدرس برگشت به سايت پذيرنده که الزاما مي بايست در دامنه سايت ثبت شده برای پذيرنده قرار داشته باشد
	payRequestData
}

type payRequestData struct {
	OrderID        uint64 `xml:"orderId"`                 // شماره درخواست(پرداخت)
	Amount         uint64 `xml:"amount"`                  // مبلغ خريد
	AdditionalData string `xml:"additionalData"`          // اطلاعات توضيحي که پذيرنده مايل به حفظ آنها برای هر تراکنش ميباشد.
	PayerID        uint64 `xml:"payerId"`                 // شناسه پرداخت کننده
	MobileNo       string `xml:"mobileNo,omitempty"`      // شماره موبايل دارنده کارت (اختياری)
	EncPan         string `xml:"encPan,omitempty"`        // شماره کارت رمز شده (اختياری)
	PanHiddenMode  string `xml:"panHiddenMode,omitempty"` // وضعیت نمايش ثابت کارت شماره (اختیاری)
	CartItem       string `xml:"cartItem,omitempty"`      // متن دلخواه پذيرنده برای نمايش در درگاه (اختیاری)
	Enc            string `xml:"enc,omitempty"`           // کارت دارنده شده رمز ملي کد ( اختیاری)
}

func transformPayRequest(in paymentgo.PayRequest) payRequestData {
	orderID, _ := strconv.ParseUint(in.OrderID, 10, 64)
	payerID, _ := strconv.ParseUint(in.UserID, 10, 64)
	return payRequestData{
		OrderID:  orderID,
		Amount:   in.Amount,
		MobileNo: in.Mobile,
		PayerID:  payerID,
	}
}

func (m service) Pay(ctx context.Context, in paymentgo.PayRequest) (string, string, error) {
	now := time.Now()
	sendDate := now.Format("20060102")
	sendTime := now.Format("150405")

	response, err := base.Call[payResponse](ctx, m.serviceURL, payRequestBody{
		BpPayRequest: payRequest{
			credentials: credentials{
				TerminalID:   m.terminalID,
				UserName:     m.username,
				UserPassword: m.password,
			},
			LocalDate:      sendDate,
			LocalTime:      sendTime,
			CallbackURL:    m.callbackURL,
			payRequestData: transformPayRequest(in),
		},
	})
	if err != nil {
		return "", "", err
	}

	parts := strings.Split(response.Result.Return, ",")
	if len(parts) < 2 || parts[0] != "0" {
		return "", "", base.ConvertError(parts[0])
	}

	refID := parts[1]
	paymentURL := fmt.Sprintf("%s?RefId=%s", m.gatewayURL, refID)
	return refID, paymentURL, nil
}

package behpardakht_test

import (
	"context"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	paymentgo "github.com/amirzayi/payment-go"
	"github.com/amirzayi/payment-go/behpardakht"
	"github.com/amirzayi/payment-go/behpardakht/base"
)

func TestVerify(t *testing.T) {
	ctx := context.Background()

	t.Run("bad request", func(t *testing.T) {
		t.Parallel()

		paymentSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		defer paymentSrv.Close()

		bp := behpardakht.NewService(http.DefaultClient, paymentSrv.URL, behpardakht.GatewayURL, "username", "password", paymentSrv.URL, 1)

		err := bp.Verify(ctx, paymentgo.VerifiyRequest{})

		var convErr *strconv.NumError
		if !errors.As(err, &convErr) {
			t.Fatalf("expected strconv.NumError, but got: %v", err)
		}
	})

	t.Run("empty response", func(t *testing.T) {
		t.Parallel()

		paymentSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		defer paymentSrv.Close()

		bp := behpardakht.NewService(http.DefaultClient, paymentSrv.URL, behpardakht.GatewayURL, "username", "password", paymentSrv.URL, 1)

		err := bp.Verify(ctx, paymentgo.VerifiyRequest{
			OrderID: "20123", ReferenceID: "48128", Amount: 5000000,
		})

		if !errors.Is(err, io.EOF) {
			t.Fatalf("expected %v, but got: %v", io.EOF, err)
		}
	})

	t.Run("non OK response http status", func(t *testing.T) {
		t.Parallel()

		paymentSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}))
		defer paymentSrv.Close()

		bp := behpardakht.NewService(http.DefaultClient, paymentSrv.URL, behpardakht.GatewayURL, "username", "password", paymentSrv.URL, 1)

		err := bp.Verify(ctx, paymentgo.VerifiyRequest{
			OrderID: "20123", ReferenceID: "48128", Amount: 5000000,
		})

		if !errors.Is(err, base.ErrInvalidResponseStatusCode) {
			t.Fatalf("expected %v, but got %v", io.EOF, err)
		}
	})

	t.Run("successful", func(t *testing.T) {
		t.Parallel()
		type (
			verifyRequest struct {
				OrderID         uint64 `xml:"orderId"`
				SaleOrderID     uint64 `xml:"saleOrderId"`
				SaleReferenceID uint64 `xml:"saleReferenceId"`
			}
			settleRequest struct {
				OrderID         uint64 `xml:"orderId"`
				SaleOrderID     uint64 `xml:"saleOrderId"`
				SaleReferenceID uint64 `xml:"saleReferenceId"`
			}
			requestBody struct {
				Body struct {
					Verify *verifyRequest `xml:"bpVerifyRequest"`
					Settle *settleRequest `xml:"bpSettleRequest"`
				} `xml:"Body"`
			}

			verifyResponse struct {
				Result base.Response `xml:"bpVerifyRequestResponse"`
			}
			settleResponse struct {
				Result base.Response `xml:"bpSettleRequestResponse"`
			}
		)

		paymentSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req requestBody
			err := xml.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if req.Body.Verify != nil {
				xml.NewEncoder(w).Encode(soapEnvelopeResponse[verifyResponse]{Body: verifyResponse{Result: base.Response{
					Return: "0",
				}}})
				return
			}
			xml.NewEncoder(w).Encode(soapEnvelopeResponse[settleResponse]{Body: settleResponse{Result: base.Response{
				Return: "0",
			}}})
		}))
		defer paymentSrv.Close()

		bp := behpardakht.NewService(http.DefaultClient, paymentSrv.URL, behpardakht.GatewayURL, "username", "password", paymentSrv.URL, 1)

		err := bp.Verify(ctx, paymentgo.VerifiyRequest{
			OrderID: "20123", ReferenceID: "48128", Amount: 5000000,
		})

		if err != nil {
			t.Fatalf("expected nil error, but got: %v", err)
		}
	})
}

type soapEnvelopeResponse[T any] struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    T        `xml:"Body"`
}

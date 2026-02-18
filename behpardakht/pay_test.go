package behpardakht_test

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	paymentgo "github.com/amirzayi/payment-go"
	"github.com/amirzayi/payment-go/behpardakht"
	"github.com/amirzayi/payment-go/behpardakht/base"
)

func TestPay(t *testing.T) {
	ctx := context.Background()

	t.Run("empty response", func(t *testing.T) {
		t.Parallel()

		paymentSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		defer paymentSrv.Close()

		bp := behpardakht.NewService(http.DefaultClient, paymentSrv.URL, behpardakht.GatewayURL, "username", "password", paymentSrv.URL, 1)

		refID, paymentURL, err := bp.Pay(ctx, paymentgo.PayRequest{
			Amount: 100000,
		})
		if !errors.Is(err, io.EOF) {
			t.Fatalf("expected %v, but got %v", io.EOF, err)
		}
		if refID != "" {
			t.Fatalf("expected empty refID but got: %s", refID)
		}
		if paymentURL != "" {
			t.Fatalf("expected empty paymentURL but got: %s", paymentURL)
		}
	})

	t.Run("non OK response http status", func(t *testing.T) {
		t.Parallel()

		paymentSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}))
		defer paymentSrv.Close()

		bp := behpardakht.NewService(http.DefaultClient, paymentSrv.URL, behpardakht.GatewayURL, "username", "password", paymentSrv.URL, 1)

		refID, paymentURL, err := bp.Pay(ctx, paymentgo.PayRequest{
			Amount: 100000,
		})
		if !errors.Is(err, base.ErrInvalidResponseStatusCode) {
			t.Fatalf("expected %v, but got %v", base.ErrInvalidResponseStatusCode, err)
		}
		if refID != "" {
			t.Fatalf("expected empty refID but got: %s", refID)
		}
		if paymentURL != "" {
			t.Fatalf("expected empty paymentURL but got: %s", paymentURL)
		}
	})

	t.Run("bad response format", func(t *testing.T) {
		t.Parallel()

		paymentSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			xml.NewEncoder(w).Encode(base.Response{
				Return: "something went wrong",
			})
			w.WriteHeader(http.StatusOK)
		}))
		defer paymentSrv.Close()

		bp := behpardakht.NewService(http.DefaultClient, paymentSrv.URL, behpardakht.GatewayURL, "username", "password", paymentSrv.URL, 1)

		refID, paymentURL, err := bp.Pay(ctx, paymentgo.PayRequest{
			Amount: 100000,
		})

		var ue xml.UnmarshalError
		if !errors.As(err, &ue) {
			t.Fatalf("expected xml.UnmarshalError, but got %v", err)
		}
		if refID != "" {
			t.Fatalf("expected empty refID but got: %s", refID)
		}
		if paymentURL != "" {
			t.Fatalf("expected empty paymentURL but got: %s", paymentURL)
		}
	})

	type payResponse struct {
		Result base.Response `xml:"bpPayRequestResponse"`
	}

	t.Run("invalid response format", func(t *testing.T) {
		t.Parallel()

		paymentSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			xml.NewEncoder(w).Encode(soapEnvelopeResponse[payResponse]{Body: payResponse{base.Response{
				Return: "something went wrong",
			}}})
			w.WriteHeader(http.StatusOK)
		}))
		defer paymentSrv.Close()

		bp := behpardakht.NewService(http.DefaultClient, paymentSrv.URL, behpardakht.GatewayURL, "username", "password", paymentSrv.URL, 1)

		refID, paymentURL, err := bp.Pay(ctx, paymentgo.PayRequest{
			Amount: 100000,
		})

		if !errors.Is(err, base.ErrInvalidResponseFormat) {
			t.Fatalf("expected %v, but got %v", base.ErrInvalidResponseFormat, err)
		}
		if refID != "" {
			t.Fatalf("expected empty refID but got: %s", refID)
		}
		if paymentURL != "" {
			t.Fatalf("expected empty paymentURL but got: %s", paymentURL)
		}
	})

	t.Run("disabled account", func(t *testing.T) {
		t.Parallel()

		paymentSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			xml.NewEncoder(w).Encode(soapEnvelopeResponse[payResponse]{Body: payResponse{base.Response{
				Return: "33",
			}}})
			w.WriteHeader(http.StatusOK)
		}))
		defer paymentSrv.Close()

		bp := behpardakht.NewService(http.DefaultClient, paymentSrv.URL, behpardakht.GatewayURL, "username", "password", paymentSrv.URL, 1)

		refID, paymentURL, err := bp.Pay(ctx, paymentgo.PayRequest{
			Amount: 100000,
		})
		if !errors.Is(err, base.ErrInvalidAccount) {
			t.Fatalf("expected %v, but got %v", base.ErrInvalidAccount, err)
		}
		if refID != "" {
			t.Fatalf("expected empty refID but got: %s", refID)
		}
		if paymentURL != "" {
			t.Fatalf("expected empty paymentURL but got: %s", paymentURL)
		}
	})

	t.Run("successful", func(t *testing.T) {
		t.Parallel()

		paymentSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			type payResponse struct {
				Result base.Response `xml:"bpPayRequestResponse"`
			}
			xml.NewEncoder(w).Encode(soapEnvelopeResponse[payResponse]{Body: payResponse{base.Response{
				Return: "0,172361",
			}}})
			w.WriteHeader(http.StatusOK)
		}))
		defer paymentSrv.Close()

		bp := behpardakht.NewService(http.DefaultClient, paymentSrv.URL, behpardakht.GatewayURL, "username", "password", paymentSrv.URL, 1)

		refID, paymentURL, err := bp.Pay(ctx, paymentgo.PayRequest{
			Amount: 100000,
		})
		if err != nil {
			t.Fatalf("expected nil error, but got: %v", err)
		}
		if refID != "172361" {
			t.Fatalf("expected 172361 refID but got: %s", refID)
		}
		expectedPaymentURL := fmt.Sprintf("%s?RefId=%s", behpardakht.GatewayURL, refID)
		if paymentURL != expectedPaymentURL {
			t.Fatalf("expected %s paymentURL but got: %s", expectedPaymentURL, paymentURL)
		}
	})
}

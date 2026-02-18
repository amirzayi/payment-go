package novinpay_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	paymentgo "github.com/amirzayi/payment-go"
	"github.com/amirzayi/payment-go/novinpay"
	"github.com/amirzayi/payment-go/novinpay/base"
	"golang.org/x/crypto/pkcs12"
)

func TestPay(t *testing.T) {
	ctx := context.Background()

	for _, tc := range []struct {
		name          string
		handler       http.HandlerFunc
		expectedError error
	}{
		{
			"empty response",
			func(w http.ResponseWriter, r *http.Request) {},
			io.EOF,
		},
		{
			"non http OK status code response",
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			base.ErrInvalidResponseStatusCode,
		},
		{
			"invalid username or password",
			func(w http.ResponseWriter, r *http.Request) {
				json.NewEncoder(w).Encode(map[string]string{
					"Result":    base.ResponseInvalidUserOrPass,
					"SessionId": "12345",
				})
			},
			base.ErrInvalidUserOrPass,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mux := http.NewServeMux()
			mux.HandleFunc("/merchantLogin/", tc.handler)
			paymentSrv := httptest.NewServer(mux)
			defer paymentSrv.Close()

			paymentService, err := novinpay.NewService(http.DefaultClient, paymentSrv.URL, novinpay.PaymentGatewayURL, "username", "password", "m123", "t123", paymentSrv.URL, "certificate_password", strings.NewReader("something"))
			if err != nil {
				t.Fatal(err)
			}
			refID, paymentURL, err := paymentService.Pay(ctx, paymentgo.PayRequest{
				Amount: 100000,
			})
			if !errors.Is(err, tc.expectedError) {
				t.Fatalf("expected %v, but got %v", tc.expectedError, err)
			}
			if refID != "" {
				t.Fatalf("expected empty refID but got: %s", refID)
			}
			if paymentURL != "" {
				t.Fatalf("expected empty paymentURL but got: %s", paymentURL)
			}
		})
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/merchantLogin/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"Result":    base.ResponseSuccess,
			"SessionId": "12345",
		})
	})
	mux.HandleFunc("/generateTransactionDataToSign/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"Result":     base.ResponseSuccess,
			"DataToSign": "12345",
			"UniqueId":   "",
		})
	})
	mux.HandleFunc("/generateSignedDataToken/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"Result":         base.ResponseSuccess,
			"ExpirationDate": time.Hour,
			"Token":          "182391",
			"ChannelId":      "",
			"UserId":         "",
		})
	})
	t.Run("invalid certificate password", func(t *testing.T) {
		t.Parallel()

		paymentSrv := httptest.NewServer(mux)
		defer paymentSrv.Close()

		f, err := os.Open("bundle.p12")
		if err != nil {
			t.Fatalf("failed to load file: %v", err)
		}
		paymentService, err := novinpay.NewService(http.DefaultClient, paymentSrv.URL, novinpay.PaymentGatewayURL, "username", "password", "m123", "t123", paymentSrv.URL, "invalid certificate password", f)
		if err != nil {
			t.Fatal(err)
		}
		refID, paymentURL, err := paymentService.Pay(ctx, paymentgo.PayRequest{
			Amount: 100000,
		})
		if !errors.Is(err, pkcs12.ErrIncorrectPassword) {
			t.Fatalf("expected %v, but got %v", pkcs12.ErrIncorrectPassword, err)
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

		paymentSrv := httptest.NewServer(mux)
		defer paymentSrv.Close()

		f, err := os.Open("bundle.p12")
		if err != nil {
			t.Fatalf("failed to load file: %v", err)
		}
		paymentService, err := novinpay.NewService(http.DefaultClient, paymentSrv.URL, novinpay.PaymentGatewayURL, "username", "password", "m123", "t123", paymentSrv.URL, "amir", f)
		if err != nil {
			t.Fatal(err)
		}
		refID, paymentURL, err := paymentService.Pay(ctx, paymentgo.PayRequest{
			Amount: 100000,
		})
		if err != nil {
			t.Fatalf("expected nil, but got %v", err)
		}
		if refID != "182391" {
			t.Fatalf("expected empty refID but got: %s", refID)
		}
		if paymentURL != novinpay.PaymentGatewayURL {
			t.Fatalf("expected empty paymentURL but got: %s", paymentURL)
		}
	})

}

package novinpay_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	paymentgo "github.com/amirzayi/payment-go"
	"github.com/amirzayi/payment-go/novinpay"
	"github.com/amirzayi/payment-go/novinpay/base"
)

func TestVerify(t *testing.T) {
	ctx := context.Background()

	for _, tc := range []struct {
		name          string
		handler       http.HandlerFunc
		expectedError error
	}{
		{
			"invald user or password",
			func(w http.ResponseWriter, r *http.Request) {
				json.NewEncoder(w).Encode(map[string]string{
					"Result": base.ResponseInvalidUserOrPass,
				})
			},
			base.ErrResponseInvalidUserOrPass,
		},
		{
			"invald source ip",
			func(w http.ResponseWriter, r *http.Request) {
				json.NewEncoder(w).Encode(map[string]string{
					"Result": base.ResponseInvalidUserOrPass,
				})
			},
			base.ErrResponseInvalidUserOrPass,
		},
		{
			"invald data",
			func(w http.ResponseWriter, r *http.Request) {
				json.NewEncoder(w).Encode(map[string]string{
					"Result": base.ResponseInvalidUserOrPass,
				})
			},
			base.ErrResponseInvalidUserOrPass,
		},
		{
			"invald order id",
			func(w http.ResponseWriter, r *http.Request) {
				json.NewEncoder(w).Encode(map[string]any{
					"Result": base.ResponseSuccess,
					"RefNum": "Some invalid ref num",
					"Amount": 512000,
				})
			},
			base.ErrMismatchVerificationRefnum,
		},
		{
			"invald amount",
			func(w http.ResponseWriter, r *http.Request) {
				json.NewEncoder(w).Encode(map[string]any{
					"Result": base.ResponseSuccess,
					"RefNum": "123",
					"Amount": 256000,
				})
			},
			base.ErrMismatchVerificationAmount,
		},
		{
			"valid",
			func(w http.ResponseWriter, r *http.Request) {
				req := make(map[string]any)
				err := json.NewDecoder(r.Body).Decode(&req)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				json.NewEncoder(w).Encode(map[string]any{
					"Result": base.ResponseSuccess,
					"RefNum": req["RefNum"],
					"Amount": 512000,
				})
			},
			nil,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/verifyMerchantTrans/", tc.handler)
			paymentSrv := httptest.NewServer(mux)
			defer paymentSrv.Close()

			paymentService, err := novinpay.NewService(paymentSrv.URL, novinpay.PaymentGatewayURL, "username", "password", "m123", "t123", paymentSrv.URL, "certificate_password", strings.NewReader("something"))
			if err != nil {
				t.Fatal(err)
			}
			err = paymentService.Verify(ctx, paymentgo.VerifiyRequest{
				OrderID:     "123",
				ReferenceID: "861",
				Amount:      512000,
			})
			if tc.expectedError == nil {
				if err != nil {
					t.Fatalf("expected nil, but got: %v", err)
				}
			} else {
				if !errors.Is(err, tc.expectedError) {
					t.Fatalf("expected %v, but got %v", tc.expectedError, err)
				}
			}
		})
	}
}

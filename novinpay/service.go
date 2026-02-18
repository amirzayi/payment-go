package novinpay

import (
	"context"
	"io"
	"net/http"
	"strings"

	paymentgo "github.com/amirzayi/payment-go"
	"github.com/amirzayi/payment-go/novinpay/base"
)

type service struct {
	httpClient          *http.Client
	serviceURL          string
	paymentGatewayURL   string
	userName            string
	password            string
	merchantID          string
	terminalID          string
	callbackURL         string
	certificatePassword string
	certificate         []byte
}

func NewService(httpClient *http.Client, serviceURL, paymentGatewayURL, userName, password, merchantID, terminalID, callbackURL, certificatePassword string, certificate io.Reader) (paymentgo.Payment, error) {
	b, err := io.ReadAll(certificate)
	if err != nil {
		return service{}, err
	}
	return service{
		httpClient:          httpClient,
		serviceURL:          serviceURL,
		paymentGatewayURL:   paymentGatewayURL,
		userName:            userName,
		password:            password,
		merchantID:          merchantID,
		terminalID:          terminalID,
		callbackURL:         callbackURL,
		certificatePassword: certificatePassword,
		certificate:         b,
	}, nil
}

type loginRequest struct {
	UserName string `json:"UserName"`
	Password string `json:"Password"`
}

type loginResponse struct {
	Result    string `json:"Result"`
	SessionId string `json:"SessionId"`
}

func (s service) login(ctx context.Context) (string, error) {
	loginResponse, err := base.DoPostApiCall[loginResponse](
		ctx,
		s.httpClient,
		s.serviceURL+loginURI,
		loginRequest{
			UserName: s.userName,
			Password: s.password,
		},
	)
	if err != nil {
		return "", err
	}
	if !strings.EqualFold(loginResponse.Result, base.ResponseSuccess) {
		return "", base.GetResponseError(loginResponse.Result)
	}
	return loginResponse.SessionId, nil
}

type wsContext struct {
	UserId   string `json:"UserId"`
	Password string `json:"Password"`
}

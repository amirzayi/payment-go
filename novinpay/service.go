package novinpay

import (
	"context"
	"os"
	"strings"

	paymentgo "github.com/amirzayi/payment-go"
	"github.com/amirzayi/payment-go/novinpay/base"
)

type service struct {
	userName            string
	password            string
	merchantID          string
	terminalID          string
	callbackURL         string
	certificatePassword string
	certificate         []byte
}

func NewService(userName, password, merchantID, terminalID, callbackURL, certificateFilePath, certificatePassword string) (paymentgo.Payment, error) {
	certificate, err := os.ReadFile(certificateFilePath)
	if err != nil {
		return service{}, err
	}
	return service{
		userName:            userName,
		password:            password,
		merchantID:          merchantID,
		terminalID:          terminalID,
		callbackURL:         callbackURL,
		certificatePassword: certificatePassword,
		certificate:         certificate,
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
	loginResponse, err := base.NewPostCall[loginResponse](
		ctx,
		loginURI,
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

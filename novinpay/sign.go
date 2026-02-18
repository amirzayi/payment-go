package novinpay

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/amirzayi/payment-go/novinpay/base"
	"golang.org/x/crypto/pkcs12"
)

type generateSignedDataTokenRequest struct {
	WsContext wsContext `json:"WSContext"`
	Signature string    `json:"Signature"`
	UniqueId  string    `json:"UniqueId"`
}

type generateSignedDataTokenResponse struct {
	Result         string        `json:"Result"`
	ExpirationDate time.Duration `json:"ExpirationDate"`
	Token          string        `json:"Token"`
	ChannelId      string        `json:"ChannelId"`
	UserId         string        `json:"UserId"`
}

func (s service) generateSignedDataToken(ctx context.Context, signature, uniqueId string) (generateSignedDataTokenResponse, error) {
	txResult, err := base.DoPostApiCall[generateSignedDataTokenResponse](
		ctx,
		s.serviceURL+generateSignedDataTokenURL,
		generateSignedDataTokenRequest{
			WsContext: wsContext{
				UserId:   s.userName,
				Password: s.password,
			},
			Signature: signature,
			UniqueId:  uniqueId,
		},
	)
	if err != nil {
		return generateSignedDataTokenResponse{}, err
	}
	if !strings.EqualFold(txResult.Result, base.ResponseSuccess) {
		return generateSignedDataTokenResponse{}, base.GetResponseError(txResult.Result)
	}
	return txResult, nil
}

func signWithKey(privateKey *rsa.PrivateKey, data string) (string, error) {
	h := crypto.SHA256.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)

	sign, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(sign), err
}

type generateTransactionRequest struct {
	WsContext              wsContext `json:"WSContext"`
	TransType              string    `json:"TransType"`
	ReserveNum             string    `json:"ReserveNum"`
	MerchantId             string    `json:"MerchantId"`
	TerminalId             string    `json:"TerminalId"`
	Amount                 uint64    `json:"Amount"`
	ProductId              string    `json:"ProductId"`
	GoodsReferenceID       string    `json:"GoodsReferenceID"`
	MerchatGoodReferenceID string    `json:"MerchatGoodReferenceID"`
	MobileNo               string    `json:"MobileNo"`
	Email                  string    `json:"Email"`
	RedirectUrl            string    `json:"RedirectUrl"`
}

type generateTransactionResponse struct {
	Result     string `json:"Result"`
	DataToSign string `json:"DataToSign"`
	UniqueId   string `json:"UniqueId"`
}

func (s service) generateTransactionDataToSign(ctx context.Context, req payRequest) (generateTransactionResponse, error) {
	txResult, err := base.DoPostApiCall[generateTransactionResponse](
		ctx,
		s.serviceURL+generateTransactionDataToSignURL,
		generateTransactionRequest{
			WsContext: wsContext{
				UserId:   s.userName,
				Password: s.password,
			},
			TransType:              "EN_GOODS",
			ReserveNum:             req.orderID,
			MerchantId:             s.merchantID,
			TerminalId:             s.terminalID,
			Amount:                 req.amount,
			ProductId:              req.orderID,
			GoodsReferenceID:       req.orderID,
			MerchatGoodReferenceID: req.orderID,
			MobileNo:               req.mobile,
			Email:                  req.email,
			RedirectUrl:            s.callbackURL,
		},
	)
	if err != nil {
		return txResult, err
	}
	if !strings.EqualFold(txResult.Result, base.ResponseSuccess) {
		return txResult, base.GetResponseError(txResult.Result)
	}
	return txResult, nil
}

func (s service) signToken(token string) (string, error) {
	privateKey, _, err := pkcs12.Decode(s.certificate, s.certificatePassword)
	if err != nil {
		return "", err
	}

	pv, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("private key is not valid")
	}
	sign, err := signWithKey(pv, token)
	if err != nil {
		return "", err
	}

	return sign, nil
}

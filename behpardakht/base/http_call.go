package base

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
)

type Response struct {
	Return string `xml:"return"`
}

type soapEnvelopeRequest struct {
	XMLName xml.Name `xml:"soapenv:Envelope"`
	SoapEnv string   `xml:"xmlns:soapenv,attr"`
	Web     string   `xml:"xmlns:web,attr"`
	Header  string   `xml:"soapenv:Header"`
	Body    any      `xml:"soapenv:Body"`
}

type soapEnvelopeResponse[T any] struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    T        `xml:"Body"`
}

func DoPostApiCall[T any](ctx context.Context, httpClient *http.Client, url string, requestBody any) (T, error) {
	var out T

	request := soapEnvelopeRequest{
		SoapEnv: "http://schemas.xmlsoap.org/soap/envelope/",
		Web:     "http://interfaces.core.sw.bps.com/",
		Body:    requestBody,
	}

	xmlData, err := xml.Marshal(request)
	if err != nil {
		return out, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(xmlData))
	if err != nil {
		return out, err
	}

	httpReq.Header.Set("Content-Type", "text/xml; charset=utf-8")

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return out, fmt.Errorf("%w %d", ErrInvalidResponseStatusCode, resp.StatusCode)
	}

	var response soapEnvelopeResponse[T]
	err = xml.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return out, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Body, nil
}

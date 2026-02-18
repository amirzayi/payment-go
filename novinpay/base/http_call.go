package base

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func DoPostApiCall[T any](ctx context.Context, httpClient *http.Client, url string, requestBody any) (T, error) {
	var response T
	body, err := json.Marshal(requestBody)
	if err != nil {
		return response, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	httpResponse, err := httpClient.Do(request)
	if err != nil {
		return response, err
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != http.StatusOK {
		return response, fmt.Errorf("%w %d", ErrInvalidResponseStatusCode, httpResponse.StatusCode)
	}

	err = json.NewDecoder(httpResponse.Body).Decode(&response)
	if err != nil {
		return response, fmt.Errorf("failed to decode response: %w", err)
	}
	return response, nil
}

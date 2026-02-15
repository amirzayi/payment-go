package base

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func DoPostApiCall[T any](ctx context.Context, url string, requestBody any) (T, error) {
	var response T
	body, err := json.Marshal(requestBody)
	if err != nil {
		return response, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	httpResponse, err := http.DefaultClient.Do(request)
	if err != nil {
		return response, err
	}
	defer httpResponse.Body.Close()

	err = json.NewDecoder(httpResponse.Body).Decode(&response)
	if err != nil {
		return response, fmt.Errorf("failed to decode response: %w", err)
	}

	if httpResponse.StatusCode != http.StatusOK {
		return response, fmt.Errorf("http call failed with status: %s", http.StatusText(httpResponse.StatusCode))
	}
	return response, nil
}

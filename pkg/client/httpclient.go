package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func CreateHttpRequest(ctx context.Context, method string, baseurl, path string, body io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", baseurl, path)
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("couldn't create request: %v", err)
	}
	return req, nil
}

func ProcessHttpResponse(resp *http.Response, expectedStatusCode int) ([]byte, int, error) {
	if resp == nil {
		return nil, 500, fmt.Errorf("nil response received")
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode != expectedStatusCode {
		return nil, resp.StatusCode, fmt.Errorf("unexpected status code: got %v; want %v", resp.StatusCode, expectedStatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("couldn't read response body: %v", err)
	}
	return body, resp.StatusCode, nil
}

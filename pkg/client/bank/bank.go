package bank

import (
	"brickbetest/internal/standarderrors"
	"brickbetest/model"
	"brickbetest/pkg/client"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	BaseUrl    string
	HTTPClient *http.Client
}

func NewClient(
	baseUrl string,
	httpClient *http.Client,
) *Client {
	return &Client{
		BaseUrl:    baseUrl,
		HTTPClient: httpClient,
	}
}

func (c Client) AccountValidation(ctx context.Context, accountNumber string) (res *GetAccountResBody, err error) {
	req, err := client.CreateHttpRequest(ctx, http.MethodGet, c.BaseUrl, string(APIEndpointGetAccount), nil)
	if err != nil {
		return nil, err
	}

	req.URL.Path = strings.Replace(req.URL.Path, ":id", accountNumber, 1)

	resp, err := c.HTTPClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("couldn't send the request: %v", err)
	}

	body, resCode, err := client.ProcessHttpResponse(resp, http.StatusOK)
	if err != nil {
		if resCode == http.StatusNotFound {
			return nil, standarderrors.NotFound
		}
		return nil, err
	}

	err = json.Unmarshal(body, &res)

	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal response body: %v", err)
	}

	return res, nil
}

func (c Client) Transfer(ctx context.Context, reqBody CreateTransferReqBody) (res *CreateTransferResBody, err error) {
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal request body: %v", err)
	}

	req, err := client.CreateHttpRequest(ctx, http.MethodPost, c.BaseUrl, string(APIEndpointCreateTransfer), bytes.NewReader(reqBodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("couldn't send the request: %v", err)
	}

	resBodyBytes, _, err := client.ProcessHttpResponse(resp, http.StatusCreated)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resBodyBytes, &res)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal response body: %v", err)
	}

	res.Status = model.TransferStatusPending // mock always return pending status

	return res, nil
}

// CheckTransferStatus mock simulate transfer status checker, 80% of it is success, 10% of it is still pending, and rest 10% will be failed
func (c Client) CheckTransferStatus(ctx context.Context, bankRefId *string) (status model.TransferStatus) {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	chance := r.Intn(100)

	time.Sleep(1 * time.Second) // mock as latency

	if chance < 80 {
		return model.TransferStatusSuccess
	} else if chance < 90 {
		return model.TransferStatusPending
	} else {
		return model.TransferStatusFailed
	}
}

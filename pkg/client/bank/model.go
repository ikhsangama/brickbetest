package bank

import "brickbetest/model"

type APIEndpoint string

const (
	APIEndpointGetAccount     APIEndpoint = "/account/:id"
	APIEndpointCreateTransfer APIEndpoint = "/transfer"
)

type GetAccountResBody struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type CreateTransferReqBody struct {
	BankCode                 string `json:"bankCode"`
	DestinationAccountNumber string `json:"destinationAccountNumber"`
	Amount                   int64  `json:"amount"`
}

type CreateTransferResBody struct {
	ReferenceId              string               `json:"referenceId"`
	DestinationAccountNumber string               `json:"destinationAccountNumber"`
	Amount                   int64                `json:"amount"`
	Status                   model.TransferStatus `json:"status"`
}

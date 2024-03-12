package model

import (
	"time"
)

type TransferReqBody struct {
	MerchantId               string `json:"merchantId"`
	Amount                   int64  `json:"amount"`
	ReferenceId              string `json:"referenceId"`
	DestinationAccountNumber string `json:"destinationAccountNumber"`
	BankCode                 string `json:"bankCode"`
}

type TransferResBody struct {
	TransferId               string         `json:"transferId"`
	MerchantId               string         `json:"merchantId"`
	MerchantRefId            string         `json:"merchantRefId"`
	Status                   TransferStatus `json:"status"`
	DestinationAccountNumber string         `json:"destinationAccountNumber"`
	Amount                   int64          `json:"amount"`
	BankCode                 string         `json:"bankCode"`
	BankRefId                string         `json:"bankRefId"`
}

type TransferStatus string

const (
	TransferStatusInit    TransferStatus = "INITIATE"
	TransferStatusPending TransferStatus = "PENDING"
	TransferStatusSuccess TransferStatus = "SUCCESS"
	TransferStatusFailed  TransferStatus = "FAILED"
)

//type Transfer struct {
//	Id                   string         `gorm:"id" json:"id"`
//	MerchantId           string         `gorm:"merchant_id" json:"merchantId"`
//	MerchantRefId        string         `gorm:"merchant_ref_id" json:"merchantRefId"`
//	BankCode             string         `gorm:"bank_code" json:"bankCode"`
//	BankRefId            *string        `gorm:"bank_ref_id" json:"bankRefId"`
//	Amount               int64          `gorm:"amount" json:"amount"`
//	Status               TransferStatus `gorm:"status" json:"status"`
//	DestinationAccNumber string         `gorm:"destination_acc_number" json:"destinationAccNumber"`
//
//	Created time.Time `gorm:"created"`
//	Updated time.Time `gorm:"updated"`
//}

type Transfer struct {
	Id                   string         `gorm:"id"`
	MerchantId           string         `gorm:"merchant_id"`
	MerchantRefId        string         `gorm:"merchant_ref_id"`
	BankCode             string         `gorm:"bank_code"`
	BankRefId            *string        `gorm:"bank_ref_id"`
	Amount               int64          `gorm:"amount"`
	Status               TransferStatus `gorm:"status"`
	DestinationAccNumber string         `gorm:"destination_acc_number"`

	Created time.Time `gorm:"created"`
	Updated time.Time `gorm:"updated"`
}

type CallbackReqBody struct {
	BankRefId string         `json:"referenceId"`
	Status    TransferStatus `json:"status"`
}

type CallbackResBody struct {
	Message string `json:"message"`
}

package model

type ValidateAccountResBody struct {
	AccountNumber string `json:"accountNumber"`
	AccountName   string `json:"accountName"`
	BankCode      string `json:"bankCode"`
}

type GetBankAccountResBody struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

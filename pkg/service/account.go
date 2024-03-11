package service

import (
	"brickbetest/model"
	"brickbetest/pkg/client/bank"
	"context"
	"log"
)

type AccountService struct {
	Client *bank.Client
}

func NewAccountService(client *bank.Client) *AccountService {
	return &AccountService{Client: client}
}

//type AccountService interface {
//	Validate(
//		ctx context.Context,
//		bankCode string,
//		accountNumber string,
//	) (resp *model.ValidateAccountResBody, err error)
//}

func (a *AccountService) Validate(ctx context.Context, bankCode string, accountNumber string) (*model.ValidateAccountResBody, error) {
	bankAccount, err := a.Client.AccountValidation(ctx, accountNumber)
	if err != nil {
		log.Printf("account validation error: %v", err)
		return nil, err
	}
	return &model.ValidateAccountResBody{
		AccountNumber: bankAccount.Id,
		AccountName:   bankAccount.Name,
		BankCode:      bankCode,
	}, nil
}

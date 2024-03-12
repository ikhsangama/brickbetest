package service

import (
	"brickbetest/internal/sqs"
	"brickbetest/internal/standarderrors"
	"brickbetest/model"
	"brickbetest/pkg/client/bank"
	"brickbetest/pkg/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type TransferService struct {
	accountService     *AccountService
	bankClient         *bank.Client
	transferRepository *repository.TransferRepository
	merchantRepository *repository.MerchantRepository
	ledgerRepository   *repository.LedgerRepository
	publisher          *sqs.Publisher
}

func NewTransferService(
	accountService *AccountService,
	bankClient *bank.Client,
	transferRepository *repository.TransferRepository,
	merchantRepository *repository.MerchantRepository,
	ledgerRepository *repository.LedgerRepository,
	publisher *sqs.Publisher,
) *TransferService {
	return &TransferService{
		accountService:     accountService,
		bankClient:         bankClient,
		transferRepository: transferRepository,
		merchantRepository: merchantRepository,
		ledgerRepository:   ledgerRepository,
		publisher:          publisher,
	}
}

func (s *TransferService) Transfer(ctx context.Context, req model.TransferReqBody) (*model.TransferResBody, error) {

	destinationBankAccountDetail, err := s.accountService.Validate(ctx, req.BankCode, req.DestinationAccountNumber)
	if err != nil {
		return nil, err
	}

	tx := s.startAndManageTransaction()
	defer func() {
		s.commitOrRollBack(tx, err)
	}()

	merchant, err := s.merchantRepository.GetAndLock(ctx, "UPDATE", req.MerchantId, tx)
	if err != nil {
		log.Printf("Failed to GetAndLock merchant with id %v, error: %v", req.MerchantId, err)
		return nil, err
	}

	newBalance := merchant.Balance - req.Amount
	if newBalance < 0 {
		log.Printf("Negative balance error for merchant with id %v, error: %v", req.MerchantId, err)
		return nil, standarderrors.InsufficientBalance
	}

	err = s.merchantRepository.UpdateBalance(ctx, req.MerchantId, newBalance, tx)
	if err != nil {
		log.Printf("Failed to UpdateBalance for merchant with id %v, error: %v", req.MerchantId, err)
		return nil, err
	}

	newTransfer, err := s.createTransfer(ctx, req, *destinationBankAccountDetail, tx)
	if err != nil {
		return nil, err
	}

	// create credit ledger
	timestamp := time.Now()
	newEntry := model.Ledger{
		Id:         uuid.New().String(),
		MerchantId: newTransfer.MerchantId,
		TransferId: newTransfer.Id,
		Credit:     newTransfer.Amount,
		Created:    timestamp,
		Updated:    timestamp,
	}
	if err := s.ledgerRepository.Create(ctx, &newEntry, tx); err != nil {
		return nil, err
	}

	if err = s.publisher.Publish(ctx, newTransfer, sqs.BankTransferRequest); err != nil {
		return nil, err
	}

	return &model.TransferResBody{
		TransferId:               newTransfer.Id,
		MerchantId:               newTransfer.MerchantId,
		MerchantRefId:            newTransfer.MerchantRefId,
		Status:                   newTransfer.Status,
		DestinationAccountNumber: newTransfer.DestinationAccNumber,
		Amount:                   newTransfer.Amount,
		BankCode:                 newTransfer.BankCode,
	}, nil
}

func (s *TransferService) GetTransfer(ctx *fasthttp.RequestCtx, id string) (*model.TransferResBody, error) {
	transfer, err := s.transferRepository.GetById(ctx, id)
	if err != nil {
		log.Printf("Error getting transfer by ID: %v", err)
		return nil, standarderrors.NotFound
	}

	return &model.TransferResBody{
		TransferId:               transfer.Id,
		MerchantId:               transfer.MerchantId,
		MerchantRefId:            transfer.MerchantRefId,
		Status:                   transfer.Status,
		DestinationAccountNumber: transfer.DestinationAccNumber,
		Amount:                   transfer.Amount,
		BankCode:                 transfer.BankCode,
	}, nil
}

func (s *TransferService) HandleTransferRequest(ctx context.Context, transfer model.Transfer) (err error) {
	tx := s.startAndManageTransaction()
	defer func() {
		s.commitOrRollBack(tx, err)
	}()

	transferReqBody := bank.CreateTransferReqBody{
		BankCode:                 transfer.BankCode,
		DestinationAccountNumber: transfer.DestinationAccNumber,
		Amount:                   transfer.Amount,
	}
	bankTransfer, err := s.bankClient.Transfer(ctx, transferReqBody)
	if err != nil {
		log.Printf("Failed to execute bank transfer, error: %v", err)
		return err
	}
	transfer.BankRefId = &bankTransfer.ReferenceId
	log.Printf("BankRefId: %v", *transfer.BankRefId)
	transfer.Status = model.TransferStatusPending
	err = s.transferRepository.Update(ctx, &transfer, tx)
	if err != nil {
		log.Printf("Failed to update transfer record, error: %v", err)
		return err
	}

	return nil
}

func (s *TransferService) TransferStatusCheck(ctx context.Context, days int, limit int) error {
	pendingTransfers, err := s.transferRepository.FindPendingTransfers(ctx, days, limit)
	if err != nil {
		return fmt.Errorf("error when finding stuck transfer requests: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(len(pendingTransfers)) // Initialize the WaitGroup with the size of pendingTransfers

	for _, tf := range pendingTransfers {
		go func(tf *model.Transfer) {
			defer wg.Done()

			status := s.bankClient.CheckTransferStatus(ctx, tf.BankRefId)
			tf.Status = status
			if tf.Status == model.TransferStatusPending {
				return
			}

			if err = s.publisher.Publish(ctx, tf, sqs.RecordTransaction); err != nil {
				log.Printf("failed to publish transaction: %v", err)
			}
		}(tf)
	}

	// Wait until all goroutines have completed
	wg.Wait()
	return nil
}

func (s *TransferService) HandleTransferCallback(ctx context.Context, req model.CallbackReqBody) (*model.CallbackResBody, error) {
	transfer, err := s.transferRepository.GetByBankRefId(ctx, req.BankRefId)
	if err != nil {
		log.Printf("failed to retrieve transaction using BankRefId %v: %v", req.BankRefId, err)
		return nil, err
	}
	transfer.Status = req.Status

	if err = s.publisher.Publish(ctx, transfer, sqs.RecordTransaction); err != nil {
		log.Printf("failed to publish transaction: %v", err)
		return nil, err
	}
	return &model.CallbackResBody{Message: "OK"}, nil
}

func (s *TransferService) CompleteTransfer(ctx context.Context, transfer model.Transfer) (err error) {
	// Database transaction init and setup
	tx := s.startAndManageTransaction()
	defer func() {
		s.commitOrRollBack(tx, err)
	}()

	// Transaction is processing
	err = s.processTransaction(ctx, transfer, tx)
	if err != nil {
		return err
	}

	return err
}

func (s *TransferService) startAndManageTransaction() *gorm.DB {
	tx := s.merchantRepository.GetDB().Begin(&sql.TxOptions{
		ReadOnly:  false,
		Isolation: sql.LevelSerializable,
	})

	if err := tx.Error; err != nil {
		log.Printf("transaction failed: %v", err)
		return nil
	}

	return tx
}

func (s *TransferService) createTransfer(
	ctx context.Context,
	req model.TransferReqBody,
	destination model.ValidateAccountResBody,
	tx *gorm.DB,
) (*model.Transfer, error) {
	timestamp := time.Now()
	newTransfer := model.Transfer{
		Id:                   uuid.New().String(),
		MerchantId:           req.MerchantId,
		Amount:               req.Amount,
		Status:               model.TransferStatusInit,
		BankCode:             req.BankCode,
		MerchantRefId:        req.ReferenceId,
		DestinationAccNumber: destination.AccountNumber,
		Created:              timestamp,
		Updated:              timestamp,
	}

	err := s.transferRepository.Create(ctx, &newTransfer, tx)
	if err != nil {
		return nil, err
	}

	return &newTransfer, nil
}

func (s *TransferService) checkIfTransferExists(ctx context.Context, req model.TransferReqBody) error {
	if err := s.checkIfTransferExists(ctx, req); err != nil {
		return nil
	}
	_, err := s.transferRepository.GetByReferenceId(ctx, req.MerchantId, req.ReferenceId)
	if err == nil {
		return standarderrors.AlreadyExist
	}
	if !errors.Is(err, standarderrors.NotFound) {
		return err
	}
	return nil
}

func (s *TransferService) processTransaction(ctx context.Context, transfer model.Transfer, tx *gorm.DB) error {
	t, err := s.transferRepository.GetAndLock(ctx, "UPDATE", transfer.Id, tx)
	if err != nil {
		if errors.Is(err, standarderrors.NotFound) {
			return standarderrors.NotFound
		}
		return err
	}

	if t.Status != model.TransferStatusPending {
		log.Println("Ignore non pending status")
		return nil
	}

	if transfer.Status == model.TransferStatusFailed {
		if err := s.processFailedTransaction(ctx, t, tx); err != nil {
			return err
		}
	}
	// else, transfer status success
	t.Status = transfer.Status
	return s.transferRepository.Update(ctx, t, tx)
}

func (s *TransferService) commitOrRollBack(tx *gorm.DB, err error) {
	if err != nil {
		rollBackErr := tx.Rollback().Error
		if rollBackErr != nil {
			log.Printf("Failed to rollback transaction: %v", rollBackErr)
		} else {
			log.Printf("Process fail, successfully rollback query transaction, err: %v", err)
		}
	} else {
		commitErr := tx.Commit().Error
		if commitErr != nil {
			log.Printf("Failed to commit transaction: %v", commitErr)
		} else {
			log.Println("Successful transaction commit")
		}
	}
}

// processFailedTransaction processes a failed transaction by creating a debit ledger and reverting the merchant's balance.
func (s *TransferService) processFailedTransaction(ctx context.Context, transfer *model.Transfer, tx *gorm.DB) error {
	m, err := s.merchantRepository.GetAndLock(ctx, "UPDATE", transfer.MerchantId, tx)
	if err != nil {
		return err
	}
	// create debit ledger
	timestamp := time.Now()
	newEntry := model.Ledger{
		Id:         uuid.New().String(),
		MerchantId: transfer.MerchantId,
		TransferId: transfer.Id,
		Debit:      transfer.Amount,
		Created:    timestamp,
		Updated:    timestamp,
	}
	if err := s.ledgerRepository.Create(ctx, &newEntry, tx); err != nil {
		return err
	}

	transfer.Status = model.TransferStatusFailed
	newBalance := m.Balance + transfer.Amount
	err = s.merchantRepository.UpdateBalance(ctx, m.Id, newBalance, tx)

	return err
}

package repository

import (
	"brickbetest/internal/standarderrors"
	"brickbetest/model"
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type TransferRepository struct {
	db *gorm.DB
}

func NewTransferRepository(db *gorm.DB) *TransferRepository {
	return &TransferRepository{
		db: db,
	}
}

func (r *TransferRepository) GetDb() *gorm.DB {
	return r.db
}

func (r *TransferRepository) GetByReferenceId(ctx context.Context, merchantId string, reference string) (result *model.Transfer, err error) {
	result = &model.Transfer{}
	if err := r.db.WithContext(ctx).Where("merchant_id=? AND merchant_ref_id=?", merchantId, reference).First(result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, standarderrors.NotFound
		}
		return nil, err
	}

	return result, nil
}

func (r *TransferRepository) Create(ctx context.Context, transfer *model.Transfer, tx *gorm.DB) error {
	executor := r.db
	if tx != nil {
		executor = tx
	}

	if err := executor.WithContext(ctx).Create(transfer).Error; err != nil {
		return err
	}

	return nil
}

func (r *TransferRepository) Update(ctx context.Context, transfer *model.Transfer, tx *gorm.DB) error {
	executor := r.db
	if tx != nil {
		executor = tx
	}

	timestamp := time.Now()
	if err := executor.WithContext(ctx).Model(transfer).Updates(map[string]any{
		"status":      transfer.Status,
		"updated":     timestamp,
		"bank_ref_id": transfer.BankRefId,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (r *TransferRepository) GetByBankRefId(ctx context.Context, id string) (result *model.Transfer, err error) {
	result = &model.Transfer{}
	if err := r.db.WithContext(ctx).Where("bank_ref_id=?", id).First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, standarderrors.NotFound
		}
		return result, err
	}

	return result, nil
}

// GetAndLock retrieves a transfer record from the database and locks it for exclusive access.
func (r *TransferRepository) GetAndLock(ctx context.Context, exclusive, id string, tx *gorm.DB) (result *model.Transfer, err error) {
	result = &model.Transfer{}
	if err := tx.WithContext(ctx).Clauses(clause.Locking{
		Strength: exclusive,
	}).Where("id=?", id).First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, standarderrors.NotFound
		}
	}

	return result, nil
}

func (r *TransferRepository) FindPendingTransfers(ctx context.Context, days int, limit int) ([]*model.Transfer, error) {
	var transfers []*model.Transfer

	conn := r.db.WithContext(ctx)

	timestamp := time.Now().AddDate(0, 0, -days)
	err := conn.Where("status = ? AND updated <= ?", model.TransferStatusPending, timestamp).Limit(limit).Find(&transfers).Error

	if err != nil {
		return nil, err
	}
	return transfers, nil
}

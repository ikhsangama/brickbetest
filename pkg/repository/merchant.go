package repository

import (
	"brickbetest/internal/standarderrors"
	"brickbetest/model"
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MerchantRepository struct {
	db *gorm.DB
}

func NewMerchantRepository(db *gorm.DB) *MerchantRepository {
	return &MerchantRepository{
		db: db,
	}
}

func (r *MerchantRepository) GetDB() *gorm.DB {
	return r.db
}

// GetAndLock retrieves a merchant record with a specified Id and locks it for exclusive access.
func (r *MerchantRepository) GetAndLock(ctx context.Context, exclusive, id string, tx *gorm.DB) (model.Merchant, error) {
	var result model.Merchant
	if err := tx.WithContext(ctx).Clauses(clause.Locking{
		Strength: exclusive,
	}).Where("id=?", id).First(&result).Error; err != nil {
		if r.isRecordNotFound(err) {
			return result, standarderrors.NotFound
		}
	}
	return result, nil
}

func (r *MerchantRepository) UpdateBalance(ctx context.Context, id string, balance int64, tx *gorm.DB) error {
	if err := tx.Model(model.Merchant{}).WithContext(ctx).Where("id=?", id).Updates(map[string]interface{}{"balance": balance}).Error; err != nil {
		if r.isRecordNotFound(err) {
			return standarderrors.NotFound
		}
	}
	return nil
}

func (r *MerchantRepository) isRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

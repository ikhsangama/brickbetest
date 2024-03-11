package repository

import (
	"brickbetest/model"
	"context"
	"gorm.io/gorm"
)

type LedgerRepository struct {
	db *gorm.DB
}

func NewLedgerRepository(db *gorm.DB) *LedgerRepository {
	return &LedgerRepository{
		db: db,
	}
}

func (l *LedgerRepository) GetDB() *gorm.DB {
	return l.db
}

func (l *LedgerRepository) Create(ctx context.Context, entry *model.Ledger, tx *gorm.DB) error {
	executor := l.db
	if tx != nil {
		executor = tx
	}

	if err := executor.WithContext(ctx).Create(entry).Error; err != nil {
		return err
	}

	return nil
}

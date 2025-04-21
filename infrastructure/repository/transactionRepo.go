package repository

import (
	"context"

	"github.com/hunick1234/phantom_mask/domain/transaction"
	"gorm.io/gorm"
)

type transactionRepoImpl struct {
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) transaction.TransactionRepo {
	return &transactionRepoImpl{db: db}
}

func (r *transactionRepoImpl) Create(t *transaction.Transaction) error {
	return r.db.Create(t).Error
}

func (r *transactionRepoImpl) Save(t *transaction.Transaction) error {
	return r.db.Save(t).Error
}
func (r *transactionRepoImpl) FindByID(id uint) (transaction.Transaction, error) {
	var t transaction.Transaction
	if err := r.db.Preload("User").Preload("Pharmacy").First(&t, id).Error; err != nil {
		return t, err
	}
	return t, nil
}

// lock 交易細節
func (r *transactionRepoImpl) WithTx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}
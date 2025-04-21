package transaction

import (
	"context"

	"gorm.io/gorm"
)

type TransactionRepo interface {
	Create(transaction *Transaction) error
	Save(transaction *Transaction) error
	FindByID(id uint) (Transaction, error)

	WithTx(ctx context.Context, fn func(tx *gorm.DB) error) error
}

package query

import (
	"time"

	"gorm.io/gorm"
)

type TransactionQueryService struct {
	db *gorm.DB
}

func NewTransactionQueryService(db *gorm.DB) *TransactionQueryService {
	return &TransactionQueryService{
		db: db,
	}
}

type TransactionSummaryQuery struct {
	StartDate time.Time `form:"start_date" time_format:"2006-01-02" binding:"required"`
	EndDate   time.Time `form:"end_date" time_format:"2006-01-02" binding:"required,gtfield=StartDate"`
}

type TransactionSummaryDTO struct {
	TotalMasks  int     `json:"total_masks"`
	TotalAmount float64 `json:"total_amount"`
}

func (t *TransactionQueryService) GetTransactionSummary(q TransactionSummaryQuery) (TransactionSummaryDTO, error) {
	var result TransactionSummaryDTO

	sql := `
	SELECT
	(SELECT SUM(transaction_amount)
	FROM transactions
	WHERE transaction_date BETWEEN $1 AND $2) AS total_amount,

	(SELECT SUM(transaction_items.quantity)
	FROM transaction_items
    JOIN transactions ON  transactions.id = transaction_items.transaction_id
	WHERE transactions.transaction_date BETWEEN $1 AND $2) AS total_masks;
	`

	err := t.db.Raw(sql, q.StartDate, q.EndDate).Scan(&result).Error
	if err != nil {
		return TransactionSummaryDTO{}, err
	}

	return result, nil
}

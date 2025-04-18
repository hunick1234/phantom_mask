package query

import "time"

type TransactionSummaryQuery struct {
	StartDate time.Time
	EndDate   time.Time
}

type TransactionSummaryDTO struct {
	TotalMasks  int     `json:"total_masks"`
	TotalAmount float64 `json:"total_amount"`
}

func (s *UserQueryService) GetTransactionSummary(q TransactionSummaryQuery) (TransactionSummaryDTO, error) {
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

	err := s.DB.Raw(sql, q.StartDate, q.EndDate).Scan(&result).Error
	if err != nil {
		return TransactionSummaryDTO{}, err
	}

	return result, nil
}

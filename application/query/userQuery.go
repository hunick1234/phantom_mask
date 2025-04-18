package query

import (
	"time"

	"gorm.io/gorm"
)

type UserQueryService struct {
	DB *gorm.DB
}

type TopUsersTransactionQuery struct {
	StartDate time.Time
	EndDate   time.Time
	Top       int
}

type TopUserDTO struct {
	UserID      string  `json:"user_id"`
	Name        string  `json:"name"`
	TotalAmount float64 `json:"total_amount"`
}


func (s *UserQueryService) GetTopUsersByTransactionAmount(q TopUsersTransactionQuery) ([]TopUserDTO, error) {
	var results []TopUserDTO

	sql := `
	SELECT users.id as user_id, users.name, SUM(transactions.transaction_amount) as total_amount
	FROM transactions
	JOIN users ON users.id = transactions.user_id
	WHERE transactions.transaction_date BETWEEN $1 AND $2
	GROUP BY users.id, users.name
	ORDER BY total_amount DESC
	LIMIT $3;
	`

	err := s.DB.Raw(sql, q.StartDate, q.EndDate,q.Top).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

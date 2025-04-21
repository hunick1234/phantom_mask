package query

import (
	"time"

	"gorm.io/gorm"
)

type UserQueryService struct {
	db *gorm.DB
}

func NewUserQuery(db *gorm.DB) *UserQueryService {
	return &UserQueryService{
		db: db,
	}
}

type TopUsersTransactionQuery struct {
	StartDate time.Time `form:"start_date" time_format:"2006-01-02" binding:"required"`
	EndDate   time.Time `form:"end_date" time_format:"2006-01-02" binding:"required,gtfield=StartDate"`
	Top       int       `form:"top" binding:"required,min=1"`
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
	WHERE transactions.transaction_date BETWEEN $1 AND $2  AND transactions.status = 'success'
	GROUP BY users.id, users.name
	ORDER BY total_amount DESC
	LIMIT $3;
	`

	err := s.db.Raw(sql, q.StartDate, q.EndDate, q.Top).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}
